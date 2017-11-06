module Action exposing (Msg(..))

import Time exposing (Time)
import Route exposing (Route)
import Page.Clients as Clients
import Page.Errored exposing (PageLoadError)


type Msg
    = ClientsLoaded (Result PageLoadError Clients.Model)
    | ClientsMsg Clients.Msg
    | LoadPage (Maybe Route)
    | SetRoute Route
    | Tick Time
