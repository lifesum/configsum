module Page.Configs exposing (Model, Msg, init, initBase, update, view)

import Date exposing (Date)
import Dict exposing (Dict)
import Html
    exposing
        ( Html
        , div
        , h1
        , h2
        , input
        , label
        , option
        , section
        , select
        , span
        , strong
        , table
        , tbody
        , td
        , text
        , th
        , thead
        , tr
        )
import Html.Attributes
    exposing
        ( checked
        , class
        , classList
        , colspan
        , disabled
        , for
        , id
        , placeholder
        , type_
        , value
        )
import Html.Events exposing (on, onCheck, onClick, onInput, targetValue)
import Http
import Json.Decode as Json
import String
import Task exposing (Task)
import Time exposing (Time)
import Api.Client
import Api.Config as Api
import Data.Client exposing (Client)
import Data.Config exposing (Config, Parameter(..))
import Page.Errored exposing (PageLoadError, pageLoadError)
import Route
import View.Error


-- MODEL


type alias Model =
    { clients : List Client
    , config : Maybe Config
    , configs : List Config
    , error : Maybe Http.Error
    , formClientId : String
    , formName : String
    , newParameter : Parameter
    , now : Time
    , showAddConfig : Bool
    , showAddParameter : Bool
    }


initModel : Time -> List Client -> Maybe Config -> List Config -> Model
initModel now clients config configs =
    Model clients config configs Nothing "" "" (BoolParameter "" False) now False False


init : Time -> Task PageLoadError Model
init now =
    let
        model clients configs =
            initModel now clients Nothing configs
    in
        Api.listBase
            |> Http.toTask
            |> Task.map2 model (Api.Client.list |> Http.toTask)
            |> Task.mapError (\err -> pageLoadError "Configs" err)


initBase : Time -> String -> Task PageLoadError Model
initBase now id =
    Api.getBase id
        |> Http.toTask
        |> Task.map (\config -> initModel now [] (Just config) [])
        |> Task.mapError (\err -> pageLoadError "Configs" err)


initParameter : String -> String -> Parameter
initParameter name typeName =
    case typeName of
        "bool" ->
            BoolParameter name False

        "number" ->
            NumberParameter name 0

        "numbers" ->
            NumbersParameter name []

        "string" ->
            StringParameter name ""

        "strings" ->
            StringsParameter name []

        _ ->
            BoolParameter name False



-- UPDATE


type Msg
    = ConfigLoaded (Result Http.Error Config)
    | FormSubmit
    | FormSubmitted (Result Http.Error Config)
    | ParameterFormSubmit Config
    | ParameterFormSubmitted (Result Http.Error Config)
    | SelectConfig String
    | ToggleAddConfig
    | ToggleAddParameter
    | UpdateFormClientId String
    | UpdateFormName String
    | UpdateParameterName String
    | UpdateParameterType String
    | UpdateValueBool Bool
    | UpdateValueNumber Int
    | UpdateValueNumbers (List Int)
    | UpdateValueString String
    | UpdateValueStrings (List String)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    let
        name =
            paramName model.newParameter
    in
        case (Debug.log "CONFIG MSG" msg) of
            ConfigLoaded (Err error) ->
                ( { model | error = Just error }, Cmd.none )

            ConfigLoaded (Ok config) ->
                ( { model | config = Just config, error = Nothing }, Cmd.none )

            FormSubmit ->
                ( { model | error = Nothing }
                , Api.createBase model.formClientId model.formName
                    |> Http.send FormSubmitted
                )

            FormSubmitted (Err error) ->
                ( model, Cmd.none )

            FormSubmitted (Ok config) ->
                ( initModel model.now model.clients Nothing (List.append model.configs [ config ]), Cmd.none )

            ParameterFormSubmit config ->
                ( model
                , Api.addParameter config.id (List.append config.parameters [ model.newParameter ])
                    |> Http.send ParameterFormSubmitted
                )

            ParameterFormSubmitted (Err error) ->
                ( { model | error = Just error }, Cmd.none )

            ParameterFormSubmitted (Ok config) ->
                ( initModel model.now model.clients (Just config) [], Cmd.none )

            SelectConfig id ->
                ( model, Route.navigate (Route.ConfigBase id) )

            ToggleAddConfig ->
                let
                    clientId =
                        case (List.head model.clients) of
                            Just client ->
                                client.id

                            Nothing ->
                                ""
                in
                    ( { model | formClientId = clientId, showAddConfig = not model.showAddConfig }, Cmd.none )

            ToggleAddParameter ->
                ( { model | showAddParameter = not model.showAddParameter }, Cmd.none )

            UpdateFormName formName ->
                ( { model | formName = formName }, Cmd.none )

            UpdateFormClientId id ->
                ( { model | formClientId = id }, Cmd.none )

            UpdateParameterName name ->
                ( { model | newParameter = paramUpdateName name model.newParameter }, Cmd.none )

            UpdateParameterType selectedType ->
                ( { model | newParameter = (initParameter name selectedType) }, Cmd.none )

            UpdateValueBool val ->
                ( { model | newParameter = (BoolParameter name val) }, Cmd.none )

            UpdateValueNumber val ->
                ( { model | newParameter = (NumberParameter name val) }, Cmd.none )

            UpdateValueNumbers vals ->
                ( { model | newParameter = (NumbersParameter name vals) }, Cmd.none )

            UpdateValueString val ->
                ( { model | newParameter = (StringParameter name val) }, Cmd.none )

            UpdateValueStrings vals ->
                ( { model | newParameter = (StringsParameter name vals) }, Cmd.none )



-- VIEW


view : Model -> Html Msg
view model =
    case model.config of
        Just config ->
            viewConfig model.now config model.showAddParameter model.newParameter model.error

        Nothing ->
            viewList model


viewAdd : Int -> String -> Msg -> Html Msg
viewAdd tdSpan lableText msg =
    tr [ class "add", onClick msg ]
        [ td [ class "type", colspan tdSpan ] [ text lableText ]
        ]


viewAddConfigForm : String -> List Client -> List (Html Msg)
viewAddConfigForm name clients =
    [ tr [ class "form" ]
        [ td []
            [ input
                [ onInput UpdateFormName
                , placeholder "Name"
                , type_ "text"
                , value name
                ]
                []
            ]
        , td []
            [ select
                [ on "change" (Json.map UpdateFormClientId targetValue) ]
                (List.map (\c -> viewOption c.name c.id) clients)
            ]
        , td []
            []
        , td []
            []
        ]
    , tr [ class "save", onClick FormSubmit ]
        [ td [ class "type", colspan 4 ] [ text "save config" ]
        ]
    ]


viewCard : ( String, String ) -> Html Msg
viewCard ( key, value ) =
    div [ class "card" ]
        [ span [] [ text key ]
        , strong [] [ text value ]
        ]


viewConfig : Time -> Config -> Bool -> Parameter -> Maybe Http.Error -> Html Msg
viewConfig now config showAdd parameter error =
    let
        action =
            if showAdd then
                viewParameterForm config parameter
            else
                [ viewAdd 3 "add parameter" ToggleAddParameter ]
    in
        div []
            [ h1 []
                [ text "Configs/Base/"
                , span [ class "highlight" ] [ text config.name ]
                ]
            , View.Error.view error
            , viewMeta config now
            , viewParameters config.parameters action
            ]


viewItem : Dict String Client -> Config -> Html Msg
viewItem clients config =
    let
        client =
            case (Dict.get config.clientId clients) of
                Just client ->
                    client.name

                Nothing ->
                    config.clientId
    in
        tr
            [ class "action"
            , (Route.href (Route.ConfigBase config.id))
            , onClick (SelectConfig config.id)
            ]
            [ td [] [ text config.name ]
            , td [] [ text client ]
            , td [] [ text config.id ]
            , td [] [ text (toString (List.length config.parameters)) ]
            ]


viewList : Model -> Html Msg
viewList model =
    let
        action =
            if model.showAddConfig then
                viewAddConfigForm model.formName model.clients
            else
                [ viewAdd 4 "add config" ToggleAddConfig
                ]

        clients =
            Dict.fromList (List.map (\c -> ( c.id, c )) model.clients)
    in
        div []
            [ h1 [] [ text "Configs/Base" ]
            , View.Error.view model.error
            , table []
                [ thead []
                    [ tr []
                        [ th [ class "name" ] [ text "name" ]
                        , th [ class "client" ] [ text "client" ]
                        , th [ class "id" ] [ text "id" ]
                        , th [] [ text "parameters" ]
                        ]
                    ]
                , tbody [] (List.append (List.map (viewItem clients) model.configs) action)
                ]
            ]


viewMeta : Config -> Time -> Html Msg
viewMeta config now =
    let
        cards =
            [ ( "id", config.id )
            , ( "client", config.clientId )
            , ( "created", (shortDate config.createdAt) )
            , ( "updated", (prettyDate now config.updatedAt) )
            ]
    in
        section [ class "meta" ] (List.map viewCard cards)


viewOption : String -> String -> Html Msg
viewOption name val =
    option [ value val ] [ text name ]


viewParameter : Parameter -> Html Msg
viewParameter parameter =
    tr []
        [ td [] [ text (paramName parameter) ]
        , td
            [ classList [ ( "type", True ), ( (paramTypeClass parameter), True ) ]
            ]
            [ text (paramTypeClass parameter)
            ]
        , td [ class ("value " ++ (paramTypeClass parameter)) ] [ viewParameterValue parameter ]
        ]


viewParameterValue : Parameter -> Html Msg
viewParameterValue parameter =
    case parameter of
        BoolParameter name value ->
            div []
                [ input
                    [ checked value
                    , disabled True
                    , id ("param-bool-" ++ name)
                    , type_ "checkbox"
                    ]
                    []
                , label [ for ("param-bool-" ++ name) ] []
                ]

        NumberParameter _ value ->
            span [] [ text (toString value) ]

        NumbersParameter _ values ->
            div [] (List.map (\v -> span [] [ text (toString v) ]) values)

        StringParameter _ value ->
            span [] [ text value ]

        StringsParameter _ values ->
            div [] (List.map (\v -> span [] [ text v ]) values)


viewParameters : List Parameter -> List (Html Msg) -> Html Msg
viewParameters parameters action =
    section [ class "parameters" ]
        [ h2 [] [ text "parameters" ]
        , table []
            [ thead []
                [ tr []
                    [ th [ class "name" ] [ text "name" ]
                    , th [ class "type" ] [ text "type" ]
                    , th [] [ text "value" ]
                    ]
                ]
            , tbody [] (List.append (List.map viewParameter parameters) action)
            ]
        ]


viewParameterForm : Config -> Parameter -> List (Html Msg)
viewParameterForm config parameter =
    let
        options =
            [ "bool"
            , "number"
            , "string"

            -- "numbers"
            -- "strings"
            ]
    in
        [ tr [ class "form" ]
            [ td [ class "name" ]
                [ input
                    [ onInput UpdateParameterName
                    , placeholder "Name"
                    , type_ "text"
                    , value (paramName parameter)
                    ]
                    []
                ]
            , td []
                [ select [ on "change" (Json.map UpdateParameterType targetValue) ] (List.map (\o -> viewOption o o) options)
                ]
            , td [ class "value" ] [ viewParameterFormValue parameter ]
            ]
        , tr [ class "save", onClick (ParameterFormSubmit config) ]
            [ td [ classList [ ( "type", True ) ], colspan 3 ] [ text "save parameter" ]
            ]
        ]


viewParameterFormValue : Parameter -> Html Msg
viewParameterFormValue parameter =
    case parameter of
        BoolParameter _ v ->
            div []
                [ input
                    [ checked v
                    , id "new-bool"
                    , onCheck UpdateValueBool
                    , type_ "checkbox"
                    ]
                    []
                , label [ for "new-bool" ] []
                ]

        NumberParameter _ v ->
            input
                [ on "input" (Json.map UpdateValueNumber targetNumber)
                , placeholder "Value"
                , type_ "text"
                , value (paramValue parameter)
                ]
                []

        NumbersParameter _ v ->
            input
                [ on "input" (Json.map UpdateValueNumbers targetNumbers)
                , placeholder "Value"
                , type_ "text"
                , value (paramValue parameter)
                ]
                []

        StringParameter _ v ->
            input
                [ onInput UpdateValueString
                , placeholder "Value"
                , type_ "text"
                , value v
                ]
                []

        StringsParameter _ v ->
            input
                [ on "input" (Json.map UpdateValueStrings targetStrings)
                , placeholder "Value"
                , type_ "text"
                , value (paramValue parameter)
                ]
                []



-- HELPER


paramName : Parameter -> String
paramName parameter =
    case parameter of
        BoolParameter name _ ->
            name

        NumberParameter name _ ->
            name

        NumbersParameter name _ ->
            name

        StringParameter name _ ->
            name

        StringsParameter name _ ->
            name


paramTypeClass : Parameter -> String
paramTypeClass parameter =
    case parameter of
        BoolParameter _ _ ->
            "bool"

        NumberParameter _ _ ->
            "number"

        NumbersParameter _ _ ->
            "numbers"

        StringParameter _ _ ->
            "string"

        StringsParameter _ _ ->
            "strings"


paramUpdateName : String -> Parameter -> Parameter
paramUpdateName name parameter =
    case parameter of
        BoolParameter _ value ->
            BoolParameter name value

        NumberParameter _ value ->
            NumberParameter name value

        NumbersParameter _ value ->
            NumbersParameter name value

        StringParameter _ value ->
            StringParameter name value

        StringsParameter _ value ->
            StringsParameter name value


paramValue : Parameter -> String
paramValue parameter =
    case parameter of
        BoolParameter _ val ->
            case val of
                True ->
                    "True"

                False ->
                    "False"

        NumberParameter _ val ->
            toString val

        NumbersParameter _ vals ->
            String.join " " (List.map (\v -> toString v) vals)

        StringParameter _ val ->
            val

        StringsParameter _ vals ->
            String.join " " vals


prettyDate : Time -> Date -> String
prettyDate now date =
    let
        day =
            24 * Time.hour

        diff =
            (now - (Date.toTime date))

        diffDays =
            Basics.floor (diff / day)
    in
        if diff < Time.minute then
            "just now"
        else if diff < (2 * Time.minute) then
            "1 minute ago"
        else if diff < Time.hour then
            (toString (Basics.floor (diff / Time.minute))) ++ " minutes ago"
        else if diff < (2 * Time.hour) then
            "1 hour ago"
        else if diff < day then
            (toString (Basics.floor (diff / Time.hour))) ++ " hours ago"
        else if diffDays == 1 then
            "yesterday"
        else if diffDays < 7 then
            (toString diffDays) ++ " days ago"
        else
            (toString (Basics.ceiling ((toFloat diffDays) / 7))) ++ " weeks ago"


shortDate : Date -> String
shortDate date =
    let
        month =
            case Date.month date of
                Date.Jan ->
                    "01"

                Date.Feb ->
                    "02"

                Date.Mar ->
                    "03"

                Date.Apr ->
                    "04"

                Date.May ->
                    "05"

                Date.Jun ->
                    "06"

                Date.Jul ->
                    "07"

                Date.Aug ->
                    "08"

                Date.Sep ->
                    "09"

                Date.Oct ->
                    "10"

                Date.Nov ->
                    "11"

                Date.Dec ->
                    "12"

        dateStr =
            String.join "-"
                [ (toString (Date.year date))
                , month
                , (toString (Date.day date))
                ]

        timeStr =
            String.join ":"
                [ (toString (Date.hour date))
                , (toString (Date.minute date))
                , (toString (Date.second date))
                ]
    in
        dateStr ++ " " ++ timeStr


toInt : String -> Int
toInt input =
    case (String.toInt input) of
        Err err ->
            0

        Ok v ->
            v


targetNumber : Json.Decoder Int
targetNumber =
    Json.map toInt (Json.at [ "target", "value" ] Json.string)


targetNumbers : Json.Decoder (List Int)
targetNumbers =
    let
        toInts input =
            List.map (\v -> toInt v) (String.split " " input)
    in
        Json.map toInts (Json.at [ "target", "value" ] Json.string)


targetStrings : Json.Decoder (List String)
targetStrings =
    let
        toStrings input =
            String.split " " input
    in
        Json.map toStrings (Json.at [ "target", "value" ] Json.string)
