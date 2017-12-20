module Data.Rule exposing (Bucket, Criteria, CriteriaUser, Kind(..), Rule, decoder)

import Date exposing (Date)
import Data.Parameter exposing (Parameter)
import Json.Decode as Decode exposing (Decoder, andThen, fail, succeed)
import Json.Decode.Pipeline exposing (decode, hardcoded, optional, required)
import Data.Parameter exposing (Parameter(..))


type Kind
    = Experiment
    | Override
    | Rollout


type alias Bucket =
    { parameters : List Parameter
    }


type alias Criteria =
    { locale : Maybe String
    , user : Maybe CriteriaUser
    }


type alias CriteriaUser =
    { id : List String
    , subscription : Maybe MatcherInt
    }


type alias MatcherInt =
    { comparator : Int
    , value : Int
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


decoder : Decoder Rule
decoder =
    decode Rule
        |> required "active" Decode.bool
        |> optional "activated_at" date (Date.fromTime 0)
        |> required "buckets" (Decode.list decodeBucket)
        |> required "config_id" Decode.string
        |> required "created_at" date
        |> optional "criteria" (Decode.map Just decodeCriteria) Nothing
        |> required "description" Decode.string
        |> optional "end_time" date (Date.fromTime 0)
        |> required "id" Decode.string
        |> required "kind" (Decode.int |> andThen decodeKind)
        |> required "name" Decode.string
        |> required "rollout" Decode.int
        |> optional "start_time" date (Date.fromTime 0)
        |> required "updated_at" date


decodeBucket : Decoder Bucket
decodeBucket =
    decode Bucket
        |> required "parameters" (Decode.list Data.Parameter.decoder)


decodeCriteria : Decoder Criteria
decodeCriteria =
    decode Criteria
        |> optional "locale" (Decode.map Just Decode.string) Nothing
        |> optional "user" (Decode.map Just decodeCriteriaUser) Nothing


decodeCriteriaUser : Decoder CriteriaUser
decodeCriteriaUser =
    decode CriteriaUser
        |> optional "id" (Decode.list Decode.string) []
        |> optional "subscription" (Decode.map Just decodeMatcherInt) Nothing


decodeKind : Int -> Decoder Kind
decodeKind raw =
    case raw of
        1 ->
            succeed Override

        2 ->
            succeed Experiment

        3 ->
            succeed Rollout

        _ ->
            fail "unsupported kind"


decodeMatcherInt : Decoder MatcherInt
decodeMatcherInt =
    decode MatcherInt
        |> required "Comparator" Decode.int
        |> required "Value" Decode.int



-- HELPER


date : Decoder Date
date =
    let
        convert : String -> Decoder Date
        convert raw =
            case Date.fromString raw of
                Ok date ->
                    succeed date

                Err error ->
                    fail error
    in
        Decode.string |> andThen convert
