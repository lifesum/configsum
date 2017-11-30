module Page.Configs exposing (Model, Msg, init, initBase, update, view)

import Date exposing (Date)
import Dict exposing (Dict)
import Html
    exposing
        ( Html
        , div
        , h1
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
        , id
        , for
        , placeholder
        , value
        , type_
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
import Data.Config exposing (Config)
import Data.Parameter exposing (Parameter(..))
import Page.Errored exposing (PageLoadError, pageLoadError)
import Route
import View.Date
import View.Error
import View.Parameter


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
            View.Parameter.name model.newParameter
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
                ( { model | error = Just error }, Cmd.none )

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
viewAdd tdSpan labelText msg =
    tr [ class "add", onClick msg ]
        [ td [ class "type", colspan tdSpan ] [ text labelText ]
        ]


viewAddConfigForm : String -> List Client -> List (Html Msg)
viewAddConfigForm name clients =
    [ tr [ class "form" ]
        [ td [ class "name" ]
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
                , strong [ class "highlight" ] [ text config.name ]
                ]
            , View.Error.view error
            , viewMeta config now
            , View.Parameter.viewTable action config.parameters
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
            [ h1 []
                [ text "Configs/"
                , strong [] [ text "Base" ]
                ]
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
            , ( "created", (View.Date.short config.createdAt) )
            , ( "updated", (View.Date.pretty now config.updatedAt) )
            ]
    in
        section [ class "meta" ] (List.map viewCard cards)


viewOption : String -> String -> Html Msg
viewOption name val =
    option [ value val ] [ text name ]


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
                    , value (View.Parameter.name parameter)
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
