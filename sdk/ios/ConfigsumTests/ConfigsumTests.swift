//
//  ConfigsumTests.swift
//  ConfigsumTests
//
//  Created by Alexandru Savu on 2017-10-04.

import XCTest
@testable import Configsum

class ConfigsumTests: XCTestCase {
    var configsum: Configsum!
    
    override func setUp() {
        super.setUp()
        let environment = Environment(log: true,
                                      token: "testToken",
                                      headers: ["X-Configsum-Userid": ["testSequence"]],
                                      baseConfigurationName: "baseConfig1",
                                      hostName: "config.svc.test3.playground.lifesum.com")
        self.configsum = Configsum(environment: environment)
    }
    
    
    override func tearDown() {
        super.tearDown()
        self.configsum = nil
    }
    
    func testGetDefaultStringValue() {
        let stringValue = self.configsum.getString(key: "stringVal1",
                                                       defaultValue: "defaultStringValue")
        XCTAssertTrue(stringValue == "defaultStringValue")
    }
    
    func testGetDefaultNumberValue() {
        let intValue = self.configsum.getInt(key: "numberVal1",
                                                 defaultValue: 1234)
        XCTAssertTrue(intValue == 1234)
    }
    
    func testGetDefaultStringListValue() {
        let stringListValue = self.configsum.getStringList(key: "stringListVal1",
                                                                defaultValue: ["defaultValue1", "defaultValue2"])
        XCTAssertTrue(stringListValue == ["defaultValue1", "defaultValue2"])
    }
    
    func testGetDefaultNumberListValue() {
        let intListValue = self.configsum.getIntList(key: "numberListVal1",
                                                         defaultValue: [1, 2, 3, 4, 5])
        XCTAssertTrue(intListValue == [1, 2, 3, 4, 5])
    }
    
    func testGetDefaultBoolValue() {
        let boolValue = self.configsum.getBool(key: "boolVal1",
                                                   defaultValue: false)
        XCTAssertFalse(boolValue)
    }
    
    func testGetRawConfig() {
        let rawConfig = self.configsum.getRawConfig()
        XCTAssertNotNil(rawConfig)
    }
}
