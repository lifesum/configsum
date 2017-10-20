//
//  HTTPClient.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-04.

import Foundation

internal struct HTTPClient: HTTP {
    internal let router: Router
    internal let logger: Logger
    
    internal init(environment: Environment) {
        self.router = Router(environment: environment)
        self.logger = Logger(environment: environment)
    }
    
    func PUT(payload: Context, completion: @escaping HTTPCompletion) {
        guard let request = router.makePUTRequest(payload) else  {
            completion(Result {
                throw HTTPError.badInput
            })
            return
        }
        execute(request: request, completion: completion)
    }
}

