//
//  Encodable.swift
//  Configsum
//
//  Created by Michel Tabari on 2017-10-17.
//  From https://github.com/zoul/generic-json-swift

import Foundation

extension Metadata: Encodable {
    
    public func encode(to encoder: Encoder) throws {
        
        var container = encoder.singleValueContainer()
        
        switch self {
        case let .array(array):
            try container.encode(array)
        case let .object(object):
            try container.encode(object)
        case let .string(string):
            try container.encode(string)
        case let .number(number):
            try container.encode(number)
        case let .bool(bool):
            try container.encode(bool)
        case .null:
            try container.encodeNil()
        }
    }
}
