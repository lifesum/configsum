module Page.Rules exposing (Model, Msg, initList, initRule, update, view)

import Date
import Html
    exposing
        ( Html
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
    Task.succeed <| Model Nothing now Nothing testList False


initRule : Time -> String -> Task PageLoadError Model
initRule now id =
    case (Decode.decodeString decoder testRulePayload) of
        Err err ->
            Task.fail <| pageLoadError "Rules" (Http.BadUrl err)

        Ok rule ->
            Task.succeed <| Model Nothing now (Just rule) [] False


testRulePayload : String
testRulePayload =
    """
    { "active": true
    , "activated_at": null
    , "buckets":
        [ { "name": "default", "parameters": [ { "name": "feature_say-cheese_toggled", "type": "bool", "value": true } ], "percentage": 0 }
        ]
    , "config_id": "01C066T0E4W2TM66RPPS6B0WN6"
    , "created_at": "2017-12-01T14:05:23.077Z"
    , "criteria": { "user": { "id": [ "123" ] } }
    , "description": "Enable say cheese for staff members."
    , "id": "01C068XFHXXRZSFHGX2A3JAB7O"
    , "kind": 1
    , "name": "override_say-cheese_staff"
    , "rollout": 0
    , "updated_at": "2017-12-01T14:05:23.077Z"
    }
    """


testList : List Rule
testList =
    [ testRule True "override" Override
    , testRule False "rollout" Rollout
    , testRule True "experiment" Experiment
    ]


testRule : Bool -> String -> Kind -> Rule
testRule active name kind =
    let
        criteria =
            Criteria <| Just (CriteriaUser testIds)

        date =
            Date.fromTime 0
    in
        Rule active date [] "config123" date Nothing "" date "id123" kind name 0 date date


testIds : List String
testIds =
    [ "17396058", "18784245", "14952160", "18636969", "14643208", "6595859", "10326163", "13818577", "17835011", "15230382", "19819697", "10116390", "14547084", "7402749", "7837787", "3920719", "10208124", "16004573", "15491054", "19651858", "12904911", "21959304", "15597571", "6097583", "18588", "11687029", "15712186", "21098618", "10326126", "24899644", "19840933", "25209715", "21231432", "8428965", "15491282", "11108767", "20456171", "10958987", "6141436", "10710556", "6807818", "6837392", "25903864", "22083683", "17963700", "19734249", "20897727", "5495849", "16925570", "7340959", "21788032", "21097143", "16756452", "19074415", "4935212", "11961267", "25547673", "5878481", "8269516", "6898344", "14544412" ]



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
            viewRule model.error model.now rule

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
