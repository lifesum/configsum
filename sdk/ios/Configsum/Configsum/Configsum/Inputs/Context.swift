//
//  Context.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-10.

import Foundation

public enum Platform: String, Codable {
    case android
    case iOS
    case watchOS
}

public struct User: Codable {
    private let age: Int?
    
    init(age: Int? = nil) {
        self.age = age
    }
}

public class Context: Codable {
    private let metadata: JSON?
    private let app: App
    private let device: Device
    private let user: User
    
    enum CodingKeys: String, CodingKey {
        case app
        case device
        case metadata
        case user
    }
    
    private struct App: Codable {
        let version: String
    }
    
    private struct OS: Codable {
        let platform: Platform
        let version: String
    }
    
    private struct Location: Codable {
        let locale: String
        let timezoneOffset: Int
    }
    
    private struct Device: Codable {
        let location: Location
        let os: OS
    }
    
    public init(appVersion: String,
                locale: String,
                platform: Platform,
                osVersion: String,
                metadata: JSON?,
                user: User) {
        let secondsOffset = TimeZone.current.secondsFromGMT()
        let location = Location(locale: locale, timezoneOffset: secondsOffset)
        let os = OS(platform: platform,
                    version: osVersion)
        self.app = App(version: appVersion)
        self.device = Device(location: location,
                            os: os)
        self.metadata = metadata
        self.user = user
    }
    
    public required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        self.app = try container.decode(App.self, forKey: .app)
        self.device = try container.decode(Device.self, forKey: .device)
        self.metadata = try container.decode(JSON.self, forKey: .metadata)
        self.user = try container.decode(User.self, forKey: .user)
    }
    
    public func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(app, forKey: .app)
        try container.encode(device, forKey: .device)
        try container.encode(metadata, forKey: .metadata)
        try container.encode(user, forKey: .user)
    }
}
