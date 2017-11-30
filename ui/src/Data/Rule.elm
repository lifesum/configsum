module Data.Rule exposing (Bucket, Criteria, CriteriaUser, Kind(..), Rule)

import Date exposing (Date)
import Data.Parameter exposing (Parameter)


type Kind
    = Experiment
    | Override
    | Rollout


type alias Bucket =
    { parameters : List Parameter
    }


type alias Criteria =
    { user : Maybe CriteriaUser
    }


type alias CriteriaUser =
    { id : List String
    }


type alias Rule =
    { active : Bool
    , activatedAt : Date
    , buckets : List Bucket
    , configId : String
    , createdAt : Date
    , criteria : Maybe Criteria
    , description : String
    , endTime : Date
    , id : String
    , kind : Kind
    , name : String
    , rollout : Int
    , startTime : Date
    , updatedAt : Date
    }
