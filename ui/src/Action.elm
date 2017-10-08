module Action exposing (Msg(..))

import Time exposing (Time)
import Route exposing (Route)


type Msg
    = LoadPage (Maybe Route)
    | SetRoute Route
    | Tick Time
