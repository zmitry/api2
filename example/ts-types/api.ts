/* eslint-disable */
// prettier-disable
// prettier-ignore
// Code generated by api2. DO NOT EDIT.
import * as t from "./types"
import {request} from "./utils"

export const api = { 
example: { 
	IEcho: {
		Hello: request<t.example.HelloRequest, t.example.HelloResponse>("POST", "/hello"),
		Echo: request<t.example.EchoRequest, t.example.EchoResponse>("POST", "/echo"),
		
	},
},
}
