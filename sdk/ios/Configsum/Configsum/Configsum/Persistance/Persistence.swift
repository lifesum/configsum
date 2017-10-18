//
//  Persistence.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-08.

import Foundation

internal struct Persistence {
    fileprivate enum Key: String {
        case values = "Configsum.config"
    }

    internal func set(result: [String: Any]) {
        UserDefaults.standard.set(result, forKey: Key.values.rawValue)
    }

    internal func get() -> [String: Any]? {
        let config = UserDefaults.standard.dictionary(forKey: Key.values.rawValue)
        return config
    }
}
