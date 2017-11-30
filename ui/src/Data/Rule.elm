module Data.Rule exposing (Rule, Kind(..))

import Date exposing (Date)


type Kind
    = Experiment
    | Override
    | Rollout


type alias Bucket =
    {}


type alias Rule =
    { active : Bool
    , activatedAt : Date
    , buckets : List Bucket
    , configId : String
    , createdAt : Date
    , description : String
    , endTime : Date
    , id : String
    , kind : Kind
    , name : String
    , rollout : Int
    , startTime : Date
    , updatedAt : Date
    }
