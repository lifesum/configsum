//
//  Context.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-10.

import Foundation

public enum Platform: String, Codable {
    case iOS
    case watchOS = "WatchOS"
}

public struct User: Codable {
    let age: Int?
    
    public init(age: Int? = nil) {
        self.age = age
    }
}

public class Context: Codable {
    private let metadata: Metadata?
    private let app: App
    private let device: Device
    private let user: User?
    private let os: OS
    private let location: Location
    
    enum CodingKeys: String, CodingKey {
        case app
        case device
        case metadata
        case user
    }
    
    private struct OS: Codable {
        let platform: Platform
        let version: String
    }
    
    private struct Location: Codable {
        let locale: String
        let timezoneOffset: Int
    }
    
    private struct App: Codable {
        let version: String
    }
    
    private struct Device: Codable {
        let location: Location
        let os: OS
    }
    
    public init(appVersion: String,
                locale: Locale,
                platform: Platform,
                osVersion: String,
                metadata: Metadata?,
                user: User?) {
        let secondsOffset = TimeZone.current.secondsFromGMT()
        
        self.location = Location(locale: locale.identifier, timezoneOffset: secondsOffset)
        self.os = OS(platform: platform,
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
        self.metadata = try container.decodeIfPresent(Metadata.self, forKey: .metadata)
        self.user = try container.decodeIfPresent(User.self, forKey: .user)
        self.os = self.device.os
        self.location = self.device.location
    }
    
    public func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(app, forKey: .app)
        try container.encode(device, forKey: .device)
        try container.encodeIfPresent(metadata, forKey: .metadata)
        try container.encodeIfPresent(user, forKey: .user)
    }
}

extension Context: CustomStringConvertible {
    public var description: String {
        let encoder = JSONEncoder()
        encoder.outputFormatting = .prettyPrinted
        let result = try! encoder.encode(self)
        return String(data: result, encoding: .utf8)!
    }
}
