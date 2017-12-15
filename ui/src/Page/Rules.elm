module Page.Rules exposing (Model, Msg, initList, initRule, update, view)

import Date
import Html
    exposing
        ( Html
        , a
        , div
        , h1
        , h2
        , section
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
import Html.Attributes exposing (class, classList, colspan)
import Html.Events exposing (onClick)
import Http
import Json.Decode as Decode
import Task exposing (Task)
import Time exposing (Time)
import Api.Rule as Api
import Data.Parameter exposing (Parameter(..))
import Data.Rule exposing (Bucket, Criteria, CriteriaUser, Kind(Experiment, Override, Rollout), Rule, decoder)
import Page.Errored exposing (PageLoadError, pageLoadError)
import View.Date
import View.Error
import View.Parameter
import Route


-- MODEL


type alias Model =
    { error : Maybe Http.Error
    , now : Time
    , rule : Maybe Rule
    , rules : List Rule
    , showAddRule : Bool
    }


initList : Time -> Task PageLoadError Model
initList now =
    Api.listRules
        |> Http.toTask
        |> Task.map (\rules -> Model Nothing now Nothing rules False)
        |> Task.mapError (\err -> pageLoadError "Rules" err)


initRule : Time -> String -> Task PageLoadError Model
initRule now id =
    Api.getRule id
        |> Http.toTask
        |> Task.map (\rule -> Model Nothing now (Just rule) [] False)
        |> Task.mapError (\err -> pageLoadError "Rules" err)



-- UPDATE


type Msg
    = ActivationToggled (Result Http.Error String)
    | FormSubmit
    | SelectRule String
    | ToggleAddRule
    | ToggleActivation String Bool


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        ActivationToggled (Err error) ->
            ( { model | error = Just error }, Cmd.none )

        ActivationToggled (Ok _) ->
            case model.rule of
                Just rule ->
                    ( { model | rule = Just { rule | active = not rule.active } }, Cmd.none )

                Nothing ->
                    ( model, Cmd.none )

        FormSubmit ->
            ( model, Cmd.none )

        SelectRule id ->
            ( model, Route.navigate <| Route.Rule id )

        ToggleAddRule ->
            ( { model | showAddRule = not model.showAddRule }, Cmd.none )

        ToggleActivation id active ->
            let
                call =
                    if active then
                        Api.deactivateRule id
                    else
                        Api.activateRule id
            in
                ( model, call |> Http.send ActivationToggled )



-- VIEW


view : Model -> Html Msg
view model =
    case model.rule of
        Just rule ->
            viewRule model.error model.now rule

        Nothing ->
            viewList model


viewActivation : Rule -> Html Msg
viewActivation rule =
    let
        ( linkClass, linkText ) =
            if rule.active then
                ( "cancel", "deactivate" )
            else
                ( "approve", "activate" )
    in
        section [ class "activation" ]
            [ h2 [] [ text "activation" ]
            , a [ class ("action " ++ linkClass), onClick <| ToggleActivation rule.id rule.active ]
                [ text linkText ]
            ]


viewAdd : Int -> String -> Msg -> Html Msg
viewAdd tdSpan labelText msg =
    tr [ class "add", onClick msg ]
        [ td [ class "type", colspan tdSpan ] [ text labelText ]
        ]


viewAddRuleForm : List (Html Msg)
viewAddRuleForm =
    [ tr [ class "form" ]
        [ td [] []
        ]
    , tr [ class "save", onClick FormSubmit ]
        [ td [ class "type", colspan 4 ] [ text "save rule" ]
        ]
    ]


viewCard : ( String, String ) -> Html Msg
viewCard ( key, value ) =
    div [ class "card" ]
        [ span [] [ text key ]
        , strong [] [ text value ]
        ]


viewCriteria : Maybe Criteria -> Html Msg
viewCriteria criteria =
    let
        attrs =
            case criteria of
                Just criteria ->
                    attrCriteriaUser criteria.user

                Nothing ->
                    []
    in
        if List.length (attrs) > 0 then
            section [ class "criteria" ]
                [ h2 [] [ text "criteria" ]
                , table []
                    [ thead []
                        [ tr []
                            [ th [ class "attribute" ] [ text "attribute" ]
                            , th [ class "match" ] [ text "match" ]
                            ]
                        ]
                    , tbody [] <| List.map viewCriteriaItem attrs
                    ]
                ]
        else
            section [ class "criteria" ] []


viewCriteriaItem : ( String, String ) -> Html Msg
viewCriteriaItem ( attr, value ) =
    tr []
        [ td [] [ text attr ]
        , td [ class "value" ] [ text value ]
        ]


viewList : Model -> Html Msg
viewList model =
    div []
        [ h1 [] [ text "Rules" ]
        , View.Error.view model.error
        , table []
            [ thead []
                [ tr []
                    [ th [ class "active icon" ] [ text "active" ]
                    , th [ class "name" ] [ text "name" ]
                    , th [ class "kind" ] [ text "kind" ]
                    , th [ class "config" ] [ text "config" ]
                    ]
                ]
            , tbody [] <| List.append (List.map viewListItem model.rules) <| viewListAction model.showAddRule
            ]
        ]


viewListAction : Bool -> List (Html Msg)
viewListAction showAddRule =
    case showAddRule of
        True ->
            viewAddRuleForm

        False ->
            [ viewAdd 4 "add rule" ToggleAddRule ]


viewListItem : Rule -> Html Msg
viewListItem rule =
    tr
        [ class "action"
        , Route.href <| Route.Rule rule.id
        , onClick <| SelectRule rule.id
        ]
        [ td [] [ text <| toString rule.active ]
        , td [] [ text rule.name ]
        , td [] [ text <| toString rule.kind ]
        , td [] [ text rule.configId ]
        ]


viewMeta : Time -> Rule -> Html Msg
viewMeta now rule =
    let
        cards =
            [ ( "id", rule.id )
            , ( "config", rule.configId )
            , ( "created", (View.Date.short rule.createdAt) )
            , ( "updated", (View.Date.pretty now rule.updatedAt) )
            , ( "activated", (View.Date.pretty now rule.activatedAt) )
            ]
    in
        section [ class "meta" ] <| List.map viewCard cards


viewRule : Maybe Http.Error -> Time -> Rule -> Html Msg
viewRule error now rule =
    div []
        [ h1 []
            [ text "Rules/"
            , strong [] [ text rule.name ]
            ]
        , View.Error.view error
        , viewMeta now rule
        , viewActivation rule
        , viewCriteria rule.criteria
        , viewParameters rule.buckets
        ]


viewParameter : Parameter -> Html Msg
viewParameter param =
    tr []
        [ td [] [ text <| View.Parameter.name param ]
        , td
            [ classList [ ( "type", True ), ( (View.Parameter.typeClass param), True ) ] ]
            [ text <| View.Parameter.typeClass param
            ]
        , td
            [ class <| "value " ++ (View.Parameter.typeClass param)
            ]
            []
        ]


viewParameters : List Bucket -> Html Msg
viewParameters buckets =
    let
        params =
            if not <| List.isEmpty buckets then
                case List.head buckets of
                    Just bucket ->
                        bucket.parameters

                    Nothing ->
                        []
            else
                []
    in
        View.Parameter.viewTable [] params



-- HELPER


attrCriteriaUser : Maybe CriteriaUser -> List ( String, String )
attrCriteriaUser user =
    case user of
        Just user ->
            [ ( "User.ID", (toString <| List.length user.id) ++ " IDs" )
            ]

        Nothing ->
            []
