//
//  Environment.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-04.

import Foundation

public struct Environment {
    public let log: Bool
    public let token: String
    public let headers: [String: [String]]
    public let baseConfigurationName: String
    public let hostName: String
    public let port: Int?
    
    /// Allows for custom environment
    ///
    /// - Parameters:
    ///   - log: flag that allows for http/https requests and responses to be printed to the console
    ///   - token: userID to be set for `X-Configsum-Userid` header
    ///   - headers: a key-value structure that holds th http headers
    ///   - baseConfigurationName: name for the base configuration
    ///   - hostName: host endpoint
    ///   - port: optional attribute for http/https requests
    public init(log: Bool,
                token: String,
                headers: [String: [String]],
                baseConfigurationName: String,
                hostName: String,
                port: Int? = nil) {
        self.log = log
        self.token = token
        self.headers = headers
        self.baseConfigurationName = baseConfigurationName
        self.hostName = hostName
        self.port = port
    }
    
    internal var serviceVersion: String {
        return "v1"
    }
}
