package self

import (
  "testing"
  "time"
)

func TestConn(t *testing.T) {
  conn, err := Init("localhost")
  if err != nil {
    t.Fail()
  }
  conn.Sub(TargetBlackboard)
  conn.Sub(TargetAgentSociety)
  conn.Sub(TargetBlackboardStream)
  conn.Sub(TargetGestureManager)
  conn.Sub(TargetSensorManager)
  time.Sleep(100 * time.Second)
  conn.Unsub(TargetBlackboard)
  conn.Unsub(TargetAgentSociety)
  conn.Unsub(TargetBlackboardStream)
  conn.Unsub(TargetGestureManager)
  conn.Unsub(TargetSensorManager)
}
