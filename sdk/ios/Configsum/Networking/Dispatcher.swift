//
//  Dispatcher.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-07.

import Foundation

internal struct Dispatcher {
    private let httpClient: HTTP
    private let persistence: Persistence
    
    internal init(httpClient: HTTP, persistence: Persistence) {
        self.httpClient = httpClient
        self.persistence = persistence
    }
    
    internal func getValue<T>(forKey key: String,
                              defaultValue: T) -> T {
        let resultHandler = ResultHandler<T>(persistence: self.persistence,
                                             key: key)
        return resultHandler.getValue(withDefault: defaultValue)
    }
    
    internal func fetchConfiguration(payload: Context,
                                     success: @escaping successBlock,
                                     failure: failureBlock?) {
        
        self.httpClient.PUT(payload: payload) { result in
            let resultHandler = ResultHandler<Any>(persistence: self.persistence)
            resultHandler.persistResult(result,
                                     success: success,
                                     failure: failure)
        }
    }
}
