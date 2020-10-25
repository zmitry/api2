
// ts-disable
// prettier-disable 
export namespace example {
export type HelloRequest = { 
	Key: string 
}

export type HelloResponse = { 
	Session: string 
}

export type EchoRequest = { 
	Session: string  
	text: string 
}

//EchoResponse
export type EchoResponse = { 
	text: string // 	Text string `json:"text"`

}

//docs
export type CustomType = { 
	Hell: number 
}

}
