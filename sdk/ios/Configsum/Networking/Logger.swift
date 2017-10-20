//
//  Logger.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-04.

import Foundation

internal struct Logger {
    internal let environment: Environment

    func log(_ request: URLRequest) {
        guard let url = request.url else { return }
        self.log("🌎 URLRequest: \(url)")
    }

    func log(_ error: Error) {
        guard let httpError = error as? HTTPError else { return }
        self.log("👎 Error Code: \(httpError.statusCode)  \(httpError.description)")
    }

    func log(_ message: String) {
        if environment.log {
            debugPrint(message)
        }
    }
}
