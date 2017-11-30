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
    public let age: Int?
    
    public init(age: Int? = nil) {
        self.age = age
    }
}

public struct OS: Codable {
    public let platform: Platform
    public let version: String
}

public struct Location: Codable {
    public let locale: String
    public let timezoneOffset: Int
}

public struct App: Codable {
    public let version: String
}

public struct Device: Codable {
    public let location: Location
    public let os: OS
}

public class Context: Codable {
    public let metadata: Metadata?
    public let app: App
    public let device: Device
    public let user: User?
    public let os: OS
    public let location: Location
    
    enum CodingKeys: String, CodingKey {
        case app
        case device
        case metadata
        case user
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
