//
//  HTTP.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-04.

import Foundation

typealias HTTPCompletion = (Result<Data>) -> Void

protocol HTTP {
    var logger: Logger { get }
    var router: Router { get }
    
    func PUT(payload: Context, completion: @escaping HTTPCompletion)
    init(environment: Environment)
}

fileprivate struct Constants {
    static let successCodes = 200...299
}

extension HTTP {
    internal func execute(request: URLRequest, completion: @escaping HTTPCompletion) {
        let session = URLSession.shared
        self.logger.log(request)
        
        let task = session.dataTask(with: request) { (data: Data?, response:URLResponse?, error:Error?) in
            completion(Result {
                if let _ = error { throw HTTPError.unknown }
                guard let data = data else {
                    self.logger.log(HTTPError.noData)
                    throw HTTPError.noData
                }
                guard let httpResponse = response as? HTTPURLResponse else {
                    self.logger.log(HTTPError.noResponse)
                    throw HTTPError.noResponse
                }
                guard self.isStatusCodeValid(httpResponse: httpResponse) else {
                    let error = HTTPError(status: httpResponse.statusCode, response: data)
                    self.logger.log(error)
                    throw error
                }
                return data
            })
        }
        task.resume()
    }
}

//Helper functions
extension HTTP {
    fileprivate func isStatusCodeValid(httpResponse: HTTPURLResponse) -> Bool {
        return Constants.successCodes.contains(httpResponse.statusCode)
    }
}
