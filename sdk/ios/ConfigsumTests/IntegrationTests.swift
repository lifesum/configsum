//
//  IntegrationTests.swift
//  ConfigsumTests
//
//  Created by Alexandru Savu on 2017-10-09.

@testable import Configsum
import XCTest

class IntegrationTests: XCTestCase {
    var configsum: Configsum!
    var attributes: Context!
    
    override func setUp() {
        super.setUp()
        let environment = Environment(log: true,
                                      token: "xCXbGXeG14GQjAPNfUKouaPJjRk68h5RAGia4wzyC1A=",
                                      headers: ["X-Configsum-Userid": ["123"]],
                                      baseConfigurationName: "houston",
                                      hostName: "localhost",
                                      port: 8700,
                                      urlScheme: "http")
        self.configsum = Configsum(environment: environment)
        self.attributes = Context(appVersion: "8.6.7",
                                     locale: Locale.current,
                                     platform: .iOS,
                                     osVersion: "8.0",
                                     metadata: nil,
                                     user: User(age: 20))
    }
    
    override func tearDown() {
        super.tearDown()
        self.configsum = nil
        self.attributes = nil
    }
    
    func testFetchConfiguration() {
        let exp = expectation(description: "testFetchConfiguration")
        self.configsum.render(attributes: self.attributes,
                              success: {
            exp.fulfill()
        }, failure: { error in
            XCTFail("http call failed with error: \(error)")
        })
        waitForExpectations(timeout: 10.0, handler: nil)
    }
    
    func testGetStringValue() {
        let exp = expectation(description: "testGetStringValue")
        self.configsum.render(attributes: attributes,
                              success: {
            let stringValue = self.configsum.getString(key: "test_string",
                                                           defaultValue: "defaultStringValue")
            XCTAssertTrue(stringValue == "houston we have a problem!")
            exp.fulfill()
        }, failure: { error in
            XCTFail("http call failed with error: \(error)")
        })
        waitForExpectations(timeout: 10.0, handler: nil)
    }
    
    func testGetNumberValue() {
        let exp = expectation(description: "test_number")
        self.configsum.render(attributes: attributes,
                              success: {
            let intValue = self.configsum.getInt(key: "test_number",
                                                     defaultValue: 1234)
            XCTAssertTrue(intValue == 42)
            exp.fulfill()
        }, failure: { error in
            XCTFail("http call failed with error: \(error)")
        })
        waitForExpectations(timeout: 10.0, handler: nil)
    }
    
    func testGetStringListValue() {
        let exp = expectation(description: "testGetStringListValue")
        self.configsum.render(attributes: attributes,
                              success: {
            let stringListValue = self.configsum.getStringList(key: "test_list_string",
                                                                   defaultValue: ["defaultValue1", "defaultValue2"])
            XCTAssertTrue(stringListValue == ["houston", "we", "have", "a", "problem"])
            exp.fulfill()
        }, failure: { error in
            XCTFail("http call failed with error: \(error)")
        })
        waitForExpectations(timeout: 10.0, handler: nil)
    }
    
    func testGetNumberListValue() {
        let exp = expectation(description: "testGetNumberListValue")
        self.configsum.render(attributes: attributes,
                              success: {
            let intListValue = self.configsum.getIntList(key: "test_list_number",
                                                             defaultValue: [1, 2, 3, 4, 5])
            XCTAssertTrue(intListValue == [2, 4, 8, 16, 32, 64, 128, 256])
            exp.fulfill()
        }, failure: { error in
            XCTFail("http call failed with error: \(error)")
        })
        waitForExpectations(timeout: 10.0, handler: nil)
    }
    
    func testGetBoolValue() {
        let exp = expectation(description: "testGetBoolValue")
        self.configsum.render(attributes: attributes,
                              success: {
            let boolValue = self.configsum.getBool(key: "test_bool", defaultValue: false)
            XCTAssertTrue(boolValue)
            exp.fulfill()
        }, failure: { error in
            XCTFail("http call failed with error: \(error)")
        })
        waitForExpectations(timeout: 10.0, handler: nil)
    }
    
    func testGetRawConfig() {
        let exp = expectation(description: "testGetRawConfig")
        self.configsum.render(attributes: self.attributes,
                              success: {
            let rawConfig = self.configsum.getRawConfig()
            XCTAssertNotNil(rawConfig)
            exp.fulfill()
        }) { error in
            XCTFail("http call failed with error: \(error)")
        }
        waitForExpectations(timeout: 10.0, handler: nil)
    }
    
    func testGetRawConfig() {
        let rawConfig = self.configsum.getRawConfig()
        XCTAssertNotNil(rawConfig)
    }
}


