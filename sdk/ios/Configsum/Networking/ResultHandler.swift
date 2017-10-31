//
//  ResultHandler.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-07.

import Foundation

enum ParseError<T, U>: Error {
    case InvalidType(have: T, want: U)
    case InvalidKey(T)
}

internal struct ResultHandler<T> {
    internal let key: String
    internal let persistence: Persistence
    
    internal init(persistence: Persistence,
                  key: String = "") {
        self.persistence = persistence
        self.key = key
    }
    
    private func toJSON(data: Data) -> Result<[String: Any]> {
        return Result{
            try JSONSerialization.jsonObject(with: data,
                                             options: .mutableContainers) as! [String: Any]
        }
    }
    
    private func value(from storedDictionary: [String: Any]) -> Any? {
        return storedDictionary[self.key]
    }
    
    private func buildDictionary(for result: Result<Data>) -> Result<[String: Any]> {
        return result.flatMap(f: self.toJSON)
    }
    
    internal func getValue(withDefault defaultValue: T) -> T {
        guard let config = persistence.get(),
            let value = self.value(from: config) else {
            return defaultValue
        }
        guard let typeValue = value as? T else {
            return defaultValue
        }
        return typeValue
    }
    
    internal func persistResult(_ result: Result<Data>,
                             success: successBlock,
                             failure: failureBlock?) {
        do {
            let configResult = self.buildDictionary(for: result)
            let config = try configResult.resolve()
            self.persistence.set(result: config)
            success()
        } catch {
            failure?(error)
        }
    }
    
    internal func rawConfig() -> [String: Any]? {
        return self.persistence.get()
    }
}
