/* eslint-disable */
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
}
export const api = { 
example: { 
	Echo: {
		Hello: request<t.example.HelloRequest, t.example.HelloResponse>("POST", "/hello"),
		Echo: request<t.example.EchoRequest, t.example.EchoResponse>("POST", "/echo"),
		
	},
},
}