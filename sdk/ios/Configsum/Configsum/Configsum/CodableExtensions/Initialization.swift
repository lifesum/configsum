//
//  Initialization.swift
//  Configsum
//
//  Created by Michel Tabari on 2017-10-17.
//  From https://github.com/zoul/generic-json-swift

import Foundation

extension Metadata {
    
    /// Create a Metadata value from anything. Argument has to be a valid Metadata structure:
    /// A `Float`, `Int`, `String`, `Bool`, an `Array` of those types or a `Dictionary`
    /// of those types.
    public init(_ value: Any) throws {
        switch value {
        case let num as Float:
            self = .number(num)
        case let num as Int:
            self = .number(Float(num))
        case let str as String:
            self = .string(str)
        case let bool as Bool:
            self = .bool(bool)
        case let array as [Any]:
            self = .array(try array.map(Metadata.init))
        case let dict as [String:Any]:
            self = .object(try dict.mapValues(Metadata.init))
        default:
            throw JSONError.decodingError
        }
    }
}

extension Metadata {
    
    /// Create a Metadata value from a `Codable`. This will give you access to the “raw”
    /// encoded Metadata value the `Codable` is serialized into. And hopefully, you could
    /// encode the resulting Metadata value and decode the original `Codable` back.
    public init<T: Codable>(codable: T) throws {
        let encoded = try JSONEncoder().encode(codable)
        self = try JSONDecoder().decode(Metadata.self, from: encoded)
    }
}

extension Metadata: ExpressibleByBooleanLiteral {
    
    public init(booleanLiteral value: Bool) {
        self = .bool(value)
    }
}

extension Metadata: ExpressibleByNilLiteral {
    
    public init(nilLiteral: ()) {
        self = .null
    }
}

extension Metadata: ExpressibleByArrayLiteral {
    
    public init(arrayLiteral elements: Metadata...) {
        self = .array(elements)
    }
}

extension Metadata: ExpressibleByDictionaryLiteral {
    
    public init(dictionaryLiteral elements: (String, Metadata)...) {
        var object: [String:Metadata] = [:]
        for (k, v) in elements {
            object[k] = v
        }
        self = .object(object)
    }
}

extension Metadata: ExpressibleByFloatLiteral {
    
    public init(floatLiteral value: Float) {
        self = .number(value)
    }
}

extension Metadata: ExpressibleByIntegerLiteral {
    
    public init(integerLiteral value: Int) {
        self = .number(Float(value))
    }
}

extension Metadata: ExpressibleByStringLiteral {
    
    public init(stringLiteral value: String) {
        self = .string(value)
    }
}
