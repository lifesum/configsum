//
//  Result.swift
//  Configsum
//
//  Created by Alexandru Savu on 2017-10-04.

enum Result<T> {
    case success(T)
    case failure(HTTPError)
    
    init(_ throwingExpr: () throws -> T) {
        do {
            let value = try throwingExpr()
            self = Result.success(value)
        } catch {
            self = Result.failure(error as! HTTPError)
        }
    }
    
    func resolve() throws -> T {
        switch self {
        case Result.success(let value): return value
        case Result.failure(let error): throw error
        }
    }
}

extension Result {
    func flatMap<U>(f: (T)->Result<U>) -> Result<U> {
        switch self {
        case .success(let t): return f(t)
        case .failure(let err): return .failure(err)
        }
    }
}
