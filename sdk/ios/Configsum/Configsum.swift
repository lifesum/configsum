//
//  Configsum.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-04.

public typealias failureBlock = (Error) -> Void
public typealias successBlock = () -> Void

public struct Configsum {
    private let dispatcher: Dispatcher
    
    
    /// Initialize a Configsum instance
    ///
    /// - Parameter environment: structure that holds parameters for building desired endpoint
    public init(environment: Environment) {
        let httpClient = HTTPClient(environment: environment)
        let persistence = Persistence()
        self.dispatcher = Dispatcher(httpClient: httpClient,
                                     persistence: persistence)
    }
    
    
    /// Returns a value from the stored config for the specified key
    ///
    /// - Parameters:
    ///   - key: a string value for accessing a rule in the current user configuration
    ///   - defaultValue: a fallback value in case of wrong key, wrong value type, or missing config
    /// - Returns: returns a boolean value associated with the rule
    public func getBool(key: String,
                        defaultValue: Bool) -> Bool {
        return dispatcher.getValue(forKey: key,
                            defaultValue: defaultValue)
    }
    
    /// Returns a value from the stored config for the specified key
    ///
    /// - Parameters:
    ///   - key: a string value for accessing a rule in the current user configuration
    ///   - defaultValue: a fallback value in case of wrong key, wrong value type, or missing config
    /// - Returns: returns a string value associated with the rule
    public func getString(key: String,
                          defaultValue: String) -> String {
        return dispatcher.getValue(forKey: key,
                            defaultValue: defaultValue)
    }
    
    /// Returns a value from the stored config for the specified key
    ///
    /// - Parameters:
    ///   - key: an Int value for accessing a rule in the current user configuration
    ///   - defaultValue: a fallback value in case of wrong key, wrong value type, or missing config
    /// - Returns: returns an Int value associated with the rule
    public func getInt(key: String,
                       defaultValue: Int) -> Int {
        return dispatcher.getValue(forKey: key,
                            defaultValue: defaultValue)
    }
    
    /// Returns a value from the stored config for the specified key
    ///
    /// - Parameters:
    ///   - key: a string list value for accessing a rule in the current user configuration
    ///   - defaultValue: a fallback value in case of wrong key, wrong value type, or missing config
    /// - Returns: returns a string list value associated with the rule
    public func getStringList(key: String,
                              defaultValue: [String]) -> [String] {
        return dispatcher.getValue(forKey: key,
                            defaultValue: defaultValue)
    }
    
    /// Returns a value from the stored config for the specified key
    ///
    /// - Parameters:
    ///   - key: an Int list value for accessing a rule in the current user configuration
    ///   - defaultValue: a fallback value in case of wrong key, wrong value type, or missing config
    /// - Returns: returns an Int list value associated with the rule
    public func getIntList(key: String,
                           defaultValue: [Int]) -> [Int] {
        return dispatcher.getValue(forKey: key,
                            defaultValue: defaultValue)
    }
    
    /// Fetch configuration through http and stores it locally
    ///
    /// - Parameters:
    ///   - attributes: a user defined structure that contains configuration data
    ///   - success: a block that executes on a successful http response
    ///   - failure: a failure block that has a HTTPError
    public func render(attributes: Context,
                       success: @escaping successBlock,
                       failure: failureBlock?) {
        dispatcher.fetchConfiguration(payload: attributes,
                                      success: success,
                                      failure: failure)
    }
    
    /// Returns the entire stored config
    ///
    /// - Returns: the config in the form for a key-value data structure
    public func getRawConfig() -> [String: Any]? {
        return dispatcher.retrieveRawConfig()
    }
}
