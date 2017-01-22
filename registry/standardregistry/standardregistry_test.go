package serviceregistry_test

import (
	"github.com/cbergoon/glb/registry"
	"github.com/cbergoon/glb/registry/standardregistry"
	"testing"
)

var sr registry.Registry = &serviceregistry.StandardRegistry{}


func init() {
	sr.Add("testSvc01", "testKey01", registry.Target{Address: "localhost:8080"})
	sr.Add("testSvc01", "testKey02", registry.Target{Address: "localhost:8081"})
	sr.Add("testSvc02", "testKey03", registry.Target{Address: "localhost"})
	sr.Add("testSvc03", "testKey04", registry.Target{Address: "localhost:80"})
	sr.Add("testSvc03", "testKey04", registry.Target{Address: "localhost:8443"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9090"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9091"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9092"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9093"})
}

func TestStandardRegistry_Lookup(t *testing.T) {
	ot, err := sr.Lookup("testSvc01", "testKey01")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 1 {
		t.Error("Expected slice length of 1 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc01", "testKey02")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 1 {
		t.Error("Expected slice length of 1 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc02", "testKey03")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 1 {
		t.Error("Expected slice length of 1 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc03", "testKey04")
	if err !=nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 2 {
		t.Error("Expected slice length of 2 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc04", "testKey05")
	if err !=nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 4 {
		t.Error("Expected slice length of 4 got ", len(ot))
	}
}

func TestStandardRegistry_Add(t *testing.T) {
	sr.Add("testSvc01", "testKey01", registry.Target{Address: "localhost:8081"})
	sr.Add("testSvc01", "testKey02", registry.Target{Address: "localhost:8082"})
	sr.Add("testSvc02", "testKey03", registry.Target{Address: "localhost"})
	sr.Add("testSvc03", "testKey04", registry.Target{Address: "localhost:81"})
	sr.Add("testSvc03", "testKey04", registry.Target{Address: "localhost:8444"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9094"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9095"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9096"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9097"})

	ot, err := sr.Lookup("testSvc01", "testKey01")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 2 {
		t.Error("Expected slice length of 2 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc01", "testKey02")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 2 {
		t.Error("Expected slice length of 2 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc02", "testKey03")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 2 {
		t.Error("Expected slice length of 2 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc03", "testKey04")
	if err !=nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 4 {
		t.Error("Expected slice length of 4 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc04", "testKey05")
	if err !=nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 8 {
		t.Error("Expected slice length of 8 got ", len(ot))
	}

	sr = &serviceregistry.StandardRegistry{}
	sr.Add("testSvc01", "testKey01", registry.Target{Address: "localhost:8080"})
	sr.Add("testSvc01", "testKey02", registry.Target{Address: "localhost:8081"})
	sr.Add("testSvc02", "testKey03", registry.Target{Address: "localhost"})
	sr.Add("testSvc03", "testKey04", registry.Target{Address: "localhost:80"})
	sr.Add("testSvc03", "testKey04", registry.Target{Address: "localhost:8443"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9090"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9091"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9092"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9093"})

}

func TestStandardRegistry_Delete(t *testing.T) {
	sr.Add("testSvc01", "testKey01", registry.Target{Address: "localhost:8081"})
	sr.Add("testSvc01", "testKey02", registry.Target{Address: "localhost:8082"})
	sr.Add("testSvc02", "testKey03", registry.Target{Address: "127.0.0.1"})
	sr.Add("testSvc03", "testKey04", registry.Target{Address: "localhost:81"})
	sr.Add("testSvc03", "testKey04", registry.Target{Address: "localhost:8444"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9094"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9095"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9096"})
	sr.Add("testSvc04", "testKey05", registry.Target{Address: "localhost:9097"})

	sr.Delete("testSvc01", "testKey01", registry.Target{Address: "localhost:8081"})
	sr.Delete("testSvc01", "testKey02", registry.Target{Address: "localhost:8082"})
	sr.Delete("testSvc02", "testKey03", registry.Target{Address: "127.0.0.1"})
	sr.Delete("testSvc03", "testKey04", registry.Target{Address: "localhost:81"})
	sr.Delete("testSvc03", "testKey04", registry.Target{Address: "localhost:8444"})
	sr.Delete("testSvc04", "testKey05", registry.Target{Address: "localhost:9094"})
	sr.Delete("testSvc04", "testKey05", registry.Target{Address: "localhost:9095"})
	sr.Delete("testSvc04", "testKey05", registry.Target{Address: "localhost:9096"})
	sr.Delete("testSvc04", "testKey05", registry.Target{Address: "localhost:9097"})

	ot, err := sr.Lookup("testSvc01", "testKey01")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 1 {
		t.Error("Expected slice length of 1 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc01", "testKey02")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 1 {
		t.Error("Expected slice length of 1 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc02", "testKey03")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 1 {
		t.Error("Expected slice length of 1 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc03", "testKey04")
	if err !=nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 2 {
		t.Error("Expected slice length of 2 got ", len(ot))
	}
	ot, err = sr.Lookup("testSvc04", "testKey05")
	if err !=nil {
		t.Error("Expected not nil error, got ", err)
	}
	if ot.Len() != 4 {
		t.Error("Expected slice length of 4 got ", len(ot))
	}
}

func TestStandardRegistry_IncrementFailures(t *testing.T) {
	counter1, err := sr.IncrementFailures("testSvc01", "testKey01", registry.Target{Address: "localhost:8080"}, 1)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter1 != 1 {
		t.Error("Expected failures value of 1 got ", counter1)
	}
	counter2, err := sr.IncrementFailures("testSvc01", "testKey01", registry.Target{Address: "localhost:8080"}, 1)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter2 != 2 {
		t.Error("Expected failures value of 2 got ", counter2)
	}
	counter3, err := sr.IncrementFailures("testSvc01", "testKey01", registry.Target{Address: "localhost:8080"}, 1)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter3 != 3 {
		t.Error("Expected failures value of 3 got ", counter3)
	}
	counter4, err := sr.IncrementFailures("testSvc02", "testKey03", registry.Target{Address: "localhost"}, 5)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter4 != 5 {
		t.Error("Expected failures value of 5 got ", counter4)
	}
	counter5, err := sr.IncrementFailures("testSvc02", "testKey03", registry.Target{Address: "localhost"}, 0)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter5 != 5 {
		t.Error("Expected failures value of 5 got ", counter5)
	}
	_, err = sr.IncrementFailures("noExist", "noExist", registry.Target{Address: "localhost:8080"}, 1)
	if err == nil {
		t.Error("Expected ErrServiceNotFound got ", err)
	}
	_, err = sr.IncrementFailures("testSvc01", "noExist", registry.Target{Address: "localhost:8080"}, 1)
	if err == nil {
		t.Error("Expected ErrServiceNotFound got ", err)
	}
}

func TestStandardRegistry_SetRoundRobbinCounter(t *testing.T) {
	counter1, err := sr.SetRoundRobbinCounter("testSvc01", "testKey01", 1)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter1 != 1 {
		t.Error("Expected counter value of 1 got ", counter1)
	}
	counter2, err := sr.SetRoundRobbinCounter("testSvc01", "testKey01", 2)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter2 != 2 {
		t.Error("Expected counter value of 2 got ", counter2)
	}
	counter3, err := sr.SetRoundRobbinCounter("testSvc01", "testKey01", 3)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter3 != 3 {
		t.Error("Expected counter value of 3 got ", counter3)
	}
	counter4, err := sr.SetRoundRobbinCounter("testSvc02", "testKey03", 5)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter4 != 5 {
		t.Error("Expected counter value of 5 got ", counter4)
	}
	counter5, err := sr.SetRoundRobbinCounter("testSvc02", "testKey03", 0)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter5 != 0 {
		t.Error("Expected counter value of 5 got ", counter5)
	}
	_, err = sr.SetRoundRobbinCounter("noExist", "noExist", 1)
	if err == nil {
		t.Error("Expected ErrServiceNotFound got ", err)
	}
	_, err = sr.SetRoundRobbinCounter("testSvc01", "noExist", 1)
	if err == nil {
		t.Error("Expected ErrServiceNotFound got ", err)
	}
}

func TestStandardRegistry_GetRoundRobbinCounter(t *testing.T) {
	counter1, err := sr.SetRoundRobbinCounter("testSvc01", "testKey01", 1)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter1 != 1 {
		t.Error("Expected counter value of 1 got ", counter1)
	}
	getCounter1, err := sr.GetRoundRobbinCounter("testSvc01", "testKey01")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if getCounter1 != 1 {
		t.Error("Expected counter value of 1 got ", getCounter1)
	}
	counter2, err := sr.SetRoundRobbinCounter("testSvc01", "testKey01", 2)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter2 != 2 {
		t.Error("Expected counter value of 2 got ", counter2)
	}
	getCounter2, err := sr.GetRoundRobbinCounter("testSvc01", "testKey01")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if getCounter2 != 2 {
		t.Error("Expected counter value of 2 got ", getCounter2)
	}
	counter3, err := sr.SetRoundRobbinCounter("testSvc01", "testKey01", 3)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter3 != 3 {
		t.Error("Expected counter value of 3 got ", counter3)
	}
	getCounter3, err := sr.GetRoundRobbinCounter("testSvc01", "testKey01")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if getCounter3 != 3 {
		t.Error("Expected counter value of 3 got ", getCounter3)
	}
	counter4, err := sr.SetRoundRobbinCounter("testSvc02", "testKey03", 5)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter4 != 5 {
		t.Error("Expected counter value of 5 got ", counter4)
	}
	getCounter4, err := sr.GetRoundRobbinCounter("testSvc02", "testKey03")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if getCounter4 != 5 {
		t.Error("Expected counter value of 5 got ", getCounter4)
	}
	counter5, err := sr.SetRoundRobbinCounter("testSvc02", "testKey03", 0)
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if counter5 != 0 {
		t.Error("Expected counter value of 5 got ", counter5)
	}
	getCounter5, err := sr.GetRoundRobbinCounter("testSvc02", "testKey03")
	if err != nil {
		t.Error("Expected not nil error, got ", err)
	}
	if getCounter5 != 0 {
		t.Error("Expected counter value of 5 got ", getCounter5)
	}
	_, err = sr.SetRoundRobbinCounter("noExist", "noExist", 1)
	if err == nil {
		t.Error("Expected ErrServiceNotFound got ", err)
	}
	_, err = sr.SetRoundRobbinCounter("testSvc01", "noExist", 1)
	if err == nil {
		t.Error("Expected ErrServiceNotFound got ", err)
	}
}
