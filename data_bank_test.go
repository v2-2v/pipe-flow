package main

import "testing"

func TestStoreDataByIndex(t *testing.T) {
    list := []Data{
        {Command: "echo \"1\"", Data: "init"},
        {Command: "echo \"2\"", Data: "init"},
        {Command: "echo \"3\"", Data: "init"},
    }

    updated := storeDataByIndex(list, 1, "2")
    if updated[1].Data != "2" {
        t.Fatalf("expected second step to receive data 2, got %q", updated[1].Data)
    }
    if updated[0].Data != "init" {
        t.Fatalf("expected first step to remain untouched, got %q", updated[0].Data)
    }

    updated = storeDataByIndex(updated, 0, "1")
    if updated[0].Data != "1" {
        t.Fatalf("expected first step to receive data 1, got %q", updated[0].Data)
    }
}
