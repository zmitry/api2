package api2

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/starius/api2/typegen"
)

type funcer interface {
	Func() interface{}
}

const tsClient = `
export const api = {
{{- range $key, $services := .}} 
{{$key}}: {
	{{- range $service, $methods := $services }} 
	{{$service}}: {
		{{range  $info := $methods -}} 
			{{$info.FnInfo.Method}}: request<t.{{.ReqType}}, t.{{.ResType}}>("{{.Method}}", "{{.Path}}"),
		{{end}}
	},{{- end}}
},{{end}}
}`

const templateHeaderDefault = `/* eslint-disable */
// prettier-disable
// prettier-ignore
import axios from "axios"
import * as t from "./types"
function cancelable<T>(p: T, source): T & { cancel: () => void } {
	let promiseAny = p as any;
	promiseAny.cancel = () => source.cancel("Request was canceled");
	let resolve = promiseAny.then.bind(p);
	promiseAny.then = (res, rej) => cancelable(resolve(res, rej), source);
	return promiseAny;
}
function request<Req, Res>(method:string, url:string) {
	return Object.assign((data: Req)=>{
		const c = axios.CancelToken.source()
		return cancelable(axios.request<Res>({ method, url, data, cancelToken: c.token  }).then(el=>el.data), c)
	}, {method, url})
}`

var tsClientTemplate = template.Must(template.New("ts_static_client").Parse(tsClient))

type TsGenConfig struct {
	OutDir         string
	TemplateHeader string
	Routes         []interface{}
	Types          []interface{}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
func GenerateTSClient(options *TsGenConfig) {
	if options.TemplateHeader == "" {
		options.TemplateHeader = templateHeaderDefault
	}
	os.RemoveAll(options.OutDir)
	os.MkdirAll(options.OutDir, os.ModePerm)
	apiFile, err := os.OpenFile(filepath.Join(options.OutDir, "api.ts"), os.O_CREATE|os.O_WRONLY, 0755)
	panicIf(err)
	_, err = apiFile.Write([]byte(options.TemplateHeader))
	panicIf(err)

	parser := typegen.NewParser()
	for _, getRoutes := range options.Routes {
		genValue := reflect.ValueOf(getRoutes)
		serviceArg := reflect.New(genValue.Type().In(0)).Elem()
		routesValues := genValue.Call([]reflect.Value{serviceArg})
		routes := routesValues[0].Interface().([]Route)
		genRoutes(apiFile, routes, parser)
	}
	typesFile, err := os.OpenFile(filepath.Join(options.OutDir, "types.ts"), os.O_WRONLY|os.O_CREATE, 0755)
	panicIf(err)
	parser.ParseRaw(options.Types...)
	typegen.PrintTsTypes(parser, typesFile)
	panicIf(err)

}

func genRoutes(w io.Writer, routes []Route, p *typegen.Parser) {
	type routeDef struct {
		Method  string
		Path    string
		ReqType interface{}
		ResType interface{}
		Handler interface{}
		FnInfo  FnInfo
	}
	m := map[string]map[string][]routeDef{}
	for _, route := range routes {
		handler := route.Handler
		if f, ok := handler.(funcer); ok {
			handler = f.Func()
		}

		handlerVal := reflect.ValueOf(handler)
		handlerType := handlerVal.Type()
		req := reflect.TypeOf(reflect.New(handlerType.In(1)).Elem().Interface()).Elem()
		response := reflect.TypeOf(reflect.New(handlerType.Out(0)).Elem().Interface()).Elem()
		p.Parse(req, response)
		fnInfo := GetFnInfo(route.Handler)
		fnInfo.StructName = strings.Replace(fnInfo.StructName, "SubService", "", 1)
		fnInfo.StructName = strings.Replace(fnInfo.StructName, "Service", "", 1)
		r := routeDef{
			ReqType: req.String(),
			ResType: response.String(),
			Method:  route.Method,
			Path:    route.Path,
			Handler: route.Handler,
			FnInfo:  fnInfo,
		}

		if _, ok := m[fnInfo.PkgName]; !ok {
			m[fnInfo.PkgName] = make(map[string][]routeDef)
		}
		m[fnInfo.PkgName][fnInfo.StructName] = append(m[fnInfo.PkgName][fnInfo.StructName], r)
	}

	err := tsClientTemplate.Execute(w, m)
	if err != nil {
		panic(err)
	}
}
