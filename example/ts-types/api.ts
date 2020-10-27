/* eslint-disable */
// prettier-disable
// prettier-ignore
// Code generated by api2. DO NOT EDIT.
import * as t from "./types"
import {request} from "./utils"

export const api = { 
example: { 
	IEchoService: { 
			Hello: request<t.example.HelloRequest, t.example.HelloResponse>("POST", "/hello", {"query":["key"]}, {"header":["session"]}), 
			Echo: request<t.example.EchoRequest, t.example.EchoResponse>("POST", "/echo", {"header":["session"],"json":["text"]}, {"json":["text"]}),
	},
},
}