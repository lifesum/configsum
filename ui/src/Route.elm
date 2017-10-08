module Route exposing (Route(..), fromLocation, href, navigate)

import Html exposing (Attribute)
import Html.Attributes
import Navigation exposing (Location, newUrl)
import UrlParser exposing (Parser, map, oneOf, parsePath, s)


type Route
    = Apps
    | Configs
    | NotFound
    | Rules


routes : Parser (Route -> a) a
routes =
    oneOf
        [ map Apps (s "")
        , map Apps (s "apps")
        , map Configs (s "configs")
        , map Rules (s "rules")
        ]


routeToString : Route -> String
routeToString route =
    let
        pieces =
            case route of
                Apps ->
                    [ "apps" ]

                Configs ->
                    [ "configs" ]

                NotFound ->
                    [ "not-found" ]

                Rules ->
                    [ "rules" ]
    in
        "/" ++ String.join "/" pieces



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
