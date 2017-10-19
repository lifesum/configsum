//
//  Decodable.swift
//  Configsum
//
//  Created by Michel Tabari on 2017-10-17.
//  From https://github.com/zoul/generic-json-swift

extension Metadata: Decodable {
    
    public init(from decoder: Decoder) throws {
        
        let container = try decoder.singleValueContainer()
        
        if let object = try? container.decode([String: Metadata].self) {
            self = .object(object)
        } else if let array = try? container.decode([Metadata].self) {
            self = .array(array)
        } else if let string = try? container.decode(String.self) {
            self = .string(string)
        } else if let bool = try? container.decode(Bool.self) {
            self = .bool(bool)
        } else if let number = try? container.decode(Float.self) {
            self = .number(number)
        } else if container.decodeNil() {
            self = .null
        } else {
            throw DecodingError.dataCorrupted(
                .init(codingPath: decoder.codingPath, debugDescription: "Invalid Metadata value.")
            )
        }
    }
}
