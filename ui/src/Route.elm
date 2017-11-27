module Route exposing (Route(..), fromLocation, href, navigate)

import Html exposing (Attribute)
import Html.Attributes
import Navigation exposing (Location, newUrl)
import UrlParser exposing (Parser, (</>), map, oneOf, parsePath, s, string)


type Route
    = Clients
    | Configs
    | ConfigsBase
    | ConfigBase String
    | NotFound
    | Rules


routes : Parser (Route -> a) a
routes =
    oneOf
        [ map Clients (s "")
        , map Clients (s "clients")
        , map Configs (s "configs")
        , map ConfigsBase (s "configs" </> s "base")
        , map ConfigBase (s "configs" </> s "base" </> string)
        , map Rules (s "rules")
        ]


routeToString : Route -> String
routeToString route =
    let
        pieces =
            case route of
                Clients ->
                    [ "clients" ]

                Configs ->
                    [ "configs" ]

                ConfigsBase ->
                    [ "configs", "base" ]

                ConfigBase id ->
                    [ "configs", "base", id ]

                NotFound ->
                    [ "not-found" ]

                Rules ->
                    [ "rules" ]
    in
        String.join "/" pieces



-- HELPERS --


fromLocation : Location -> Maybe Route
fromLocation location =
    parsePath routes location


href : Route -> Attribute msg
href route =
    Html.Attributes.href (routeToString route)


navigate : Route -> Cmd msg
navigate route =
    newUrl (routeToString route)
