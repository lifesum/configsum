//
//  ContextTests.swift
//  ConfigsumTests
//
//  Created by Michel Tabari on 2017-10-17.
//

import XCTest
@testable import Configsum

class ContextTests: XCTestCase {
    
    var encoder: JSONEncoder!
    var decoder: JSONDecoder!
    var contextWithoutMetadata = "contextWithoutMetadata"
    var contextWithMetadata = "contextWithMetadata"
    var contextWithComplexMetadata = "contextWithComplexMetadata"
    
    override func setUp() {
        super.setUp()
        self.encoder = JSONEncoder()
        self.encoder.outputFormatting = .prettyPrinted
        self.decoder = JSONDecoder()
    }
    
    override func tearDown() {
        self.encoder = nil
        self.decoder = nil
    }
    
    func testEncodeContextWithoutMetadata() {
        // Read the Metadata data from the test file
        let jsonFilePath = Bundle(for: type(of: self)).path(forResource: self.contextWithoutMetadata, ofType: "json")
        let dataFromJSONFile = NSData(contentsOfFile: jsonFilePath!)!
        
        let wantedContext = try! decoder.decode(Context.self,
                                                from: dataFromJSONFile as Data)
        
        let context = Context(appVersion: "8.6.0",
                              locale: Locale.current,
                              platform: .iOS,
                              osVersion: "11.0",
                              metadata: nil,
                              user: User(age: 20))
        
        let resultString = String(data: try! encoder.encode(context), encoding: .utf8)!
        let expectedString = String(data: try! encoder.encode(wantedContext), encoding: .utf8)!
        
        XCTAssertEqual(resultString, expectedString)
    }
    
    func testEncodeContextWithMetadata() {
        // Read the Metadata data from the test file
        let jsonFilePath = Bundle(for: type(of: self)).path(forResource: self.contextWithMetadata, ofType: "json")
        let dataFromJSONFile = NSData(contentsOfFile: jsonFilePath!)!
        
        let wantedContext = try! decoder.decode(Context.self,
                                                from: dataFromJSONFile as Data)
        
        let context = Context(appVersion: "8.6.0",
                              locale: Locale.current,
                              platform: .iOS,
                              osVersion: "11.0",
                              metadata: ["name": "testName",
                                         "age": 22],
                              user: User(age: 20))
        
        let resultString = String(data: try! encoder.encode(context), encoding: .utf8)!
        let expectedString = String(data: try! encoder.encode(wantedContext), encoding: .utf8)!
        
        XCTAssertEqual(resultString, expectedString)
    }
    
    func testEncodeContextWithComplexMetadata() {
        // Read the Metadata data from the test file
        let jsonFilePath = Bundle(for: type(of: self)).path(forResource: self.contextWithComplexMetadata, ofType: "json")
        let dataFromJSONFile = NSData(contentsOfFile: jsonFilePath!)!
        
        let wantedContext = try! decoder.decode(Context.self,
                                                from: dataFromJSONFile as Data)
        
        let context = Context(appVersion: "8.6.0",
                              locale: Locale.current,
                              platform: .iOS,
                              osVersion: "11.0",
                              metadata: ["name": "testName",
                                         "age": 22,
                                         "nestedDictionary": ["nestedStringList": ["item1", "item2"],
                                                              "nestedBool": true]],
                              user: User(age: 20))
        
        let resultString = String(data: try! encoder.encode(context), encoding: .utf8)!
        let expectedString = String(data: try! encoder.encode(wantedContext), encoding: .utf8)!
        
        XCTAssertEqual(resultString, expectedString)
    }
}
