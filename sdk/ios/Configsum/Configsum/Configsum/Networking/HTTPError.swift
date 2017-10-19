//
//  HTTPError.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-04.

import Foundation

enum HTTPError: Error {
    case noData
    case noResponse
    case badInput
    case notFound(Data)
    case badRequest(Data)
    case unauthorized(Data)
    case internalServerError(Data)
    case serviceUnavailable(Data)
    case gatewayTimeout
    case unknown
    
    init(status: Int, response: Data) {
        switch status {
        case 400: self = .badRequest(response)
        case 401: self = .unauthorized(response)
        case 404: self = .notFound(response)
        case 500: self = .internalServerError(response)
        case 503: self = .serviceUnavailable(response)
        case 504: self = .gatewayTimeout
        default: self = .unknown
        }
    }
    
    var statusCode: Int {
        switch self {
        case .badRequest: return 400
        case .unauthorized: return 401
        case .notFound: return 404
        case .internalServerError: return 500
        case .serviceUnavailable: return 503
        case .gatewayTimeout: return 504
        default: return -1
        }
    }
    
    var description: String {
        switch self {
        case .noData: return "HTTP error: no data"
        case .noResponse: return "HTTP error: no response object"
        case .badInput: return "Bad JSON input"
        case let .notFound(body): return parseErrorResponse(body)
        case let .badRequest(body): return parseErrorResponse(body)
        case let .unauthorized(body): return parseErrorResponse(body)
        case let .internalServerError(body): return parseErrorResponse(body)
        case let .serviceUnavailable(body): return parseErrorResponse(body)
        case .gatewayTimeout: return "HTTP error: Gateway Timeout"
        default: return "unknown"
        }
    }
    
    private func parseErrorResponse(_ data: Data) -> String {
        do {
            let responseJSON = try JSONSerialization.jsonObject(with: data,
                                             options: .mutableContainers) as! [String: Any]
            guard let errorBody = responseJSON["reason"] as? String else {
                return "Error body has wrong format"
            }
            return errorBody
        } catch {
            return "Could not parse error body"
        }
    }
}

