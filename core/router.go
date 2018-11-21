package core

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func NewRouter(init func(r *Router)) interface{} {
	return &SimpleRouter{Init: init}
}

func NewCoreRouter(bindMux ...*mux.Router) *Router {
	defaultRouter := mux.NewRouter()
	if len(bindMux) > 0 {
		defaultRouter = bindMux[0]
	}
	defaultRouter.MethodNotAllowedHandler = http.HandlerFunc(SendResponse(NewErrorResponse(HTTPErrors.Get(http.StatusMethodNotAllowed))))
	defaultRouter.NotFoundHandler = http.HandlerFunc(SendResponse(NewErrorResponse(HTTPErrors.Get(http.StatusNotFound))))

	return &Router{
		mux: defaultRouter,
		Errors: NewErrorMap(),
		RouterConfig: &RouterConfig{
			guardConfig: &GuardConfigMap{},
		},
	}
}

func (r *Router) GetHandler() *mux.Router {
	return r.mux
}

func (r *Router) Get(path string, resolver Resolver) IRouterConfig {
	return r.Handle(path, resolver, http.MethodGet)
}

func (r *Router) Put(path string, resolver Resolver) IRouterConfig {
	return r.Handle(path, resolver, http.MethodPut)
}

func (r *Router) Post(path string, resolver Resolver) IRouterConfig {
	return r.Handle(path, resolver, http.MethodPost)
}

func (r *Router) Delete(path string, resolver Resolver) IRouterConfig {
	return r.Handle(path, resolver, http.MethodDelete)
}

func (r *Router) All(path string, resolver Resolver) IRouterConfig {
	return r.Handle(path, resolver)
}

func (r *Router) Handle(path string, resolver Resolver, methods ...string) IRouterConfig {
	handlerFunc := r.mux.HandleFunc(path, r.resolver(resolver))
	if len(methods) > 0 {
		handlerFunc.Methods(methods...)
	}
	return r
}

func (r *Router) Doc(explanation string) IRouterConfig {
	return r
}

func (r *Router) DocPath(path string) IRouterConfig {
	return r
}

func (r *Router) Guard(config interface{}, passIfError ...bool) IRouterConfig {
	if reflect.ValueOf(config).Kind() != reflect.Struct {
		RouterLogger.Err("Cannot use config that is not a struct")
		return nil
	}

	r.guard = config

	defaultPassIfError := false
	if len(passIfError) > 0 {
		defaultPassIfError = passIfError[0]
	}

	r.passIfError = defaultPassIfError

	return r
}

func (r *Router) HandleWebSocket(path string, ws *WebSocket) {
	r.mux.HandleFunc(path, ws.handler)
}

func (r *Router) Use(middlewareResolver MiddlewareResolver) {
	r.mux.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			profiler := TrackTime()
			response := middlewareResolver(r.getMiddlewareContext(res, req, handler, profiler))
			if response == (Response{}) {
				return
			}
			if response.Status == http.StatusContinue {
				handler.ServeHTTP(res, req)
			} else {
				SendResponse(response)(res, req)
			}
		})
	})
}

func (r *Router) Set(path string, router interface{}) {
	routerPrefix := r.mux.PathPrefix(path).Subrouter()
	newPrefixedRouter := NewCoreRouter(routerPrefix)

	routerValue := reflect.ValueOf(router).Elem()
	field := routerValue.FieldByName("Router")

	if field.IsValid() {
		field.Set(reflect.ValueOf(newPrefixedRouter))
	} else {
		return
	}

	routerType := reflect.TypeOf(router)
	method, ok := routerType.MethodByName(router_init_func_name)
	if !ok {
		simpleMethod := routerValue.FieldByName(router_init_func_name)

		if simpleMethod != reflect.Zero(reflect.TypeOf(simpleMethod)).Interface() {
			simpleMethod.Call([]reflect.Value{reflect.ValueOf(newPrefixedRouter)})
		} else {
			RouterLogger.Warn("Missing Init method in router!")
		}
		return
	}

	method.Func.Call([]reflect.Value{reflect.ValueOf(router)})
}

func (r *Router) SetHandler(path string, handler http.Handler) {
	r.mux.Handle(path, handler)
}

func (r *Router) resolver(resolver Resolver) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		profiler := TrackTime()

		if r.guard != nil {
			body, err := r.validateGuard(req)
			RouterLogger.Log(err)
			if err != nil {
				if r.passIfError {
					r.guard = err
				} else {
					SendResponse(NewErrorResponse(NewError(http.StatusBadRequest, err.Error())))(res, req)
					return
				}
			}
			r.guard = body
		}

		response := resolver(r.getContext(res, req, profiler))

		if response != (Response{}) {
			SendResponse(response)(res, req)
		}
	}
}

func (r *Router) validateGuard(req *http.Request) (interface{}, error) {
	r.guardConfig = r.newGuardConfigMap(r.guard)

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &r.guard)
	if err != nil {
		return nil, err
	}

	guardBody, err := r.validateGuardConfig(r.guardConfig, r.guard)
	if err != nil {
		return nil, err
	}

	return guardBody, nil
}

func (r *Router) validateGuardConfig(config *GuardConfigMap, body interface{}) (interface{}, error) {
	bodyType := reflect.TypeOf(body)

	RouterLogger.Info(config)

	if bodyType.Kind() != reflect.Map {
		return nil, errors.New("body is not a map")
	}

	bodyMap := body.(map[string]interface{})

	for key, field := range bodyMap {
		fieldType := reflect.TypeOf(field).Kind()
		configField, ok := r.getConfigByFieldName(config, key)
		if !ok {
			return nil, errors.New("guard with name " + configField.FieldName + " doesn't exist")
		}

		RouterLogger.Info(configField)

		if fieldType == reflect.Map {
			RouterLogger.Info("ITS A MAP", key)
			_, err := r.validateGuardConfig(configField.Node, field)
			if err != nil {
				return nil, err
			}
		}

		if fieldType == reflect.String {
			valueString := field.(string)

			if len(strings.Trim(valueString, " ")) == 0 && !configField.Omitempty {
				return nil, errors.New("field " + configField.FieldName + " is empty")
			}

			if val := configField.Options[option_contains]; val != nil {
				if len(val.([]string)) > 0 && strings.Trim(val.([]string)[0], " ") != "" {
					for _, n := range val.([]string) {
						if !strings.Contains(valueString, n) {
							return nil, errors.New("the " + configField.FieldName + " filed does not contain a string/char: " + n)
						}
					}
				}
			}

			if val := configField.Options[option_deprecate]; val != nil {
				if len(val.([]string)) > 0 && strings.Trim(val.([]string)[0], " ") != "" {
					for _, n := range val.([]string) {
						if strings.Contains(valueString, n) {
							return nil, errors.New("the " + configField.FieldName + " filed does contains a forbidden string/char: " + n)
						}
					}
				}
			}
		}

		if fieldType != reflect.Bool && fieldType != reflect.Struct {
			RouterLogger.Info(fieldType, key)
			if fieldType == reflect.Slice {
				valueSlice := field.([]interface{})

				if val := configField.Options[option_max]; val != nil && len(valueSlice) > int(val.(int)) {
					return nil, errors.New("the " + configField.FieldName + " field has exceeded the max value")
				}

				if val := configField.Options[option_min]; val != nil && len(valueSlice) < int(val.(int)) {
					return nil, errors.New("the " + configField.FieldName + " field has exceeded the min value")
				}
			} else if fieldType == reflect.String {
				valueString := field.(string)

				if val := configField.Options[option_max]; val != nil && len(valueString) > int(val.(int)) {
					return nil, errors.New("the " + configField.FieldName + " field has exceeded the max value")
				}

				if val := configField.Options[option_min]; val != nil && len(valueString) < int(val.(int)) {
					return nil, errors.New("the " + configField.FieldName + " field has exceeded the min value")
				}
			} else {
				valueInt := int64(field.(int))

				if val := configField.Options[option_max]; val != nil && valueInt > int64(val.(int)) {
					return nil, errors.New("the " + configField.FieldName + " field has exceeded the max value")
				}
				if val := configField.Options[option_min]; val != nil && valueInt < int64(val.(int)) {
					return nil, errors.New("the " + configField.FieldName + " field has exceeded the min value")
				}
			}
		}
	}

	return body, nil
}

func (r *Router) getConfigByFieldName(config *GuardConfigMap, field string) (*GuardConfig, bool) {
	for _, n := range (*config) {
		RouterLogger.Info(n.FieldName, "==", field)
		if n.FieldName == field {
			RouterLogger.Info(true)
			return n, true
		}
	}
	RouterLogger.Warn("Stop here")

	return &GuardConfig{}, false
}

func (r *Router) newGuardConfigMap(guard interface{}) *GuardConfigMap {
	guardConfigMap := &GuardConfigMap{}
	guardType := reflect.TypeOf(guard)
	guardValue := reflect.ValueOf(guard)

	for i := 0; i < guardType.NumField(); i++ {
		field := guardType.Field(i)

		guardConfig := &GuardConfig{
			Options: map[string]interface{}{},
		}

		jsonTag := field.Tag.Get("json")

		guardConfig.FieldName = r.getFieldNameFromJSONTag(strings.Replace(jsonTag, " ", "", len(jsonTag)), field.Name)
		guardConfig.Type = field.Type.Kind()

		untrimmedTagString := field.Tag.Get(winter_guard_tag)
		tagString := strings.Replace(untrimmedTagString, " ", "", len(untrimmedTagString))
		if len(tagString) == 0 {
			continue
		}

		for _, n := range strings.Split(tagString, ",") {
			if len(n) == 0 {
				continue
			}

			option := strings.Split(n, ":")
			if len(option) != 2 {
				switch option[0] {
				case option_unrequired:
					guardConfig.Omitempty = true
				default:
					RouterLogger.Warn("Unknown guard tag option:", option[0])
				}
				continue
			}

			optionKey := option[0]
			optionVal := option[1]

			checkLength := guardConfig.Type != reflect.Bool && guardConfig.Type != reflect.Struct

			switch optionKey {
			case option_extends:
				parent, ok := (*guardConfigMap)[optionVal]
				if !ok {
					RouterLogger.Warn("Field " + optionVal + " doesn't exists in GuardMap, check the extends option and order of fields")
					continue
				}
				guardConfig.Omitempty = parent.Omitempty
				guardConfig.Options[option_max] = parent.Options[option_max]
				guardConfig.Options[option_min] = parent.Options[option_min]
				guardConfig.Options[option_contains] = parent.Options[option_contains]
				guardConfig.Options[option_deprecate] = parent.Options[option_deprecate]
			case option_max:
				if checkLength {
					maxLen, err := strconv.Atoi(optionVal)
					if err != nil {
						RouterLogger.Warn("Value of max option isn't looks like a number")
						continue
					}
					guardConfig.Options[option_max] = maxLen
				} else {
					RouterLogger.Warn("Max option is not usable for field with type", guardConfig.Type.String())
				}
			case option_min:
				if checkLength {
					minLen, err := strconv.Atoi(optionVal)
					if err != nil {
						RouterLogger.Warn("Value of min option isn't looks like a number")
						continue
					}
					guardConfig.Options[option_min] = minLen
				} else {
					RouterLogger.Warn("Min option is not usable for field with type", guardConfig.Type.String())
				}
			case option_contains:
				if guardConfig.Type == reflect.String {
					guardConfig.Options[option_contains] = r.getCharsFromPrefixedString(optionVal)
				} else {
					RouterLogger.Warn("Contains option is not usable for field with type", guardConfig.Type.String())
				}
			case option_deprecate:
				if guardConfig.Type == reflect.String {
					guardConfig.Options[option_deprecate] = r.getCharsFromPrefixedString(optionVal)
				} else {
					RouterLogger.Warn("Deprecated option is not usable for field with type", guardConfig.Type.String())
				}
			default:
				RouterLogger.Warn("Unknown guard tag option:", optionKey)
				continue
			}
		}

		if field.Type.Kind() == reflect.Struct {
			guardConfig.Node = r.newGuardConfigMap(guardValue.Field(i).Interface())
		}

		(*guardConfigMap)[field.Name] = guardConfig
	}

	return guardConfigMap
}

func (r *Router) getCharsFromPrefixedString(v string) []string {
	charArr := strings.Split(v, "'~'")
	chars := []string{}
	charArrLen := len(charArr) - 1

	if len(charArr) == 0 {
		return []string{}
	} else if len(charArr) == 1 {
		return []string{charArr[0][1:len(charArr[0]) - 1]}
	}

	chars = append(chars, (charArr[0])[1:])

	for i := 1; i < len(charArr) - 1; i++ {
		chars = append(chars, charArr[i])
	}

	chars = append(chars, (charArr[charArrLen])[:len(charArr[charArrLen]) - 1])

	return chars
}

func (r *Router) getFieldNameFromJSONTag(tag, field string) string {
	options := strings.Split(tag, ",")
	for _, n := range options {
		if n == "omitempty" {
			continue
		}
		return n
	}
	return field
}

func (r *Router) getContext(res http.ResponseWriter, req *http.Request, executionTracker func() time.Duration) *Context {
	return &Context{
		guard: r.guard,
		Request: req,
		Response: res,
		TrackTime: executionTracker,
	}
}

func (r *Router) getMiddlewareContext(res http.ResponseWriter, req *http.Request, handler http.Handler, executionTracker func() time.Duration) *MiddlewareContext {
	return &MiddlewareContext{
		Context: r.getContext(res, req, executionTracker),
		handler: handler,
	}
}
