//
//  Router.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-04.

import Foundation

internal struct Router {
    fileprivate struct Constants {
        static let tokenHeader = "X-Configsum-Token"
        static let authorizationKey = "Authorization"
        static let contentType = "application/json"
        static let contentTypeKey = "Content-Type"
    }
    
    let environment: Environment
    private var requestURL: URL? {
        var components = URLComponents()
        components.scheme = environment.urlScheme
        components.host = environment.hostName
        components.port = environment.port
        components.path = "/\(environment.serviceVersion)/config/\(environment.baseConfigurationName)"
        return components.url
    }

    internal func makePUTRequest(_ payload: Context) -> URLRequest? {
        var request = URLRequest(url: requestURL!)
        guard let encodedPayload = encodePayload(payload) else { return nil }
        request.httpMethod = "PUT"
        request.addValue(environment.token, forHTTPHeaderField: Constants.tokenHeader)
        request.addValue(Constants.contentType, forHTTPHeaderField: Constants.contentTypeKey)
        request.httpBody = encodedPayload
        for (field, values) in environment.headers {
            values.forEach{ request.addValue($0, forHTTPHeaderField: field) }
        }
        return request
    }
    
    private func encodePayload(_ payload: Context) -> Data? {
        do {
            let encoder = JSONEncoder()
            return try encoder.encode(payload)
        } catch {
            return nil
        }
    }
}
