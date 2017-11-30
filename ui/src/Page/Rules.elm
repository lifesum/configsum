module Page.Rules exposing (Model, Msg, initList, initRule, update, view)

import Date
import Html exposing (Html, div, h1, table, tbody, td, text, th, thead, tr, strong)
import Html.Attributes exposing (class, colspan)
import Html.Events exposing (onClick)
import Http
import Task exposing (Task)
import Data.Rule exposing (Rule, Kind(Experiment, Override, Rollout))
import Page.Errored exposing (PageLoadError)
import View.Error
import Route


-- MODEL


type alias Model =
    { error : Maybe Http.Error
    , rule : Maybe Rule
    , rules : List Rule
    , showAddRule : Bool
    }


initList : Task PageLoadError Model
initList =
    Task.succeed <| Model Nothing Nothing testList False


initRule : String -> Task PageLoadError Model
initRule id =
    Task.succeed <| Model Nothing (Just (testRule True "override" Override)) [] False


testList : List Rule
testList =
    [ testRule True "override" Override
    , testRule False "rollout" Rollout
    , testRule True "experiment" Experiment
    ]


testRule : Bool -> String -> Kind -> Rule
testRule active name kind =
    let
        date =
            Date.fromTime 0
    in
        Rule active date [] "config123" date "" date "id123" kind name 0 date date



-- UPDATE


type Msg
    = FormSubmit
    | SelectRule String
    | ToggleAddRule


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FormSubmit ->
            ( model, Cmd.none )

        SelectRule id ->
            ( model, Route.navigate <| Route.Rule id )

        ToggleAddRule ->
            ( { model | showAddRule = not model.showAddRule }, Cmd.none )



-- VIEW


view : Model -> Html Msg
view model =
    case model.rule of
        Just rule ->
            viewRule rule

        Nothing ->
            viewList model


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


viewRule : Rule -> Html Msg
viewRule rule =
    div []
        [ h1 []
            [ text "Rules/"
            , strong [] [ text rule.name ]
            ]
        ]
