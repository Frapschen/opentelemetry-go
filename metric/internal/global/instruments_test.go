// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package global

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/noop"
)

func testFloat64Race(interact func(float64), setDelegate func(metric.Meter)) {
	finish := make(chan struct{})
	go func() {
		for {
			interact(1)
			select {
			case <-finish:
				return
			default:
			}
		}
	}()

	setDelegate(noop.NewMeterProvider().Meter(""))
	close(finish)
}

func testInt64Race(interact func(int64), setDelegate func(metric.Meter)) {
	finish := make(chan struct{})
	go func() {
		for {
			interact(1)
			select {
			case <-finish:
				return
			default:
			}
		}
	}()

	setDelegate(noop.NewMeterProvider().Meter(""))
	close(finish)
}

func TestAsyncInstrumentSetDelegateRace(t *testing.T) {
	// Float64 Instruments
	t.Run("Float64", func(t *testing.T) {
		t.Run("Counter", func(t *testing.T) {
			delegate := &afCounter{}
			f := func(float64) { _ = delegate.Unwrap() }
			testFloat64Race(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(t *testing.T) {
			delegate := &afUpDownCounter{}
			f := func(float64) { _ = delegate.Unwrap() }
			testFloat64Race(f, delegate.setDelegate)
		})

		t.Run("Gauge", func(t *testing.T) {
			delegate := &afGauge{}
			f := func(float64) { _ = delegate.Unwrap() }
			testFloat64Race(f, delegate.setDelegate)
		})
	})

	// Int64 Instruments

	t.Run("Int64", func(t *testing.T) {
		t.Run("Counter", func(t *testing.T) {
			delegate := &aiCounter{}
			f := func(int64) { _ = delegate.Unwrap() }
			testInt64Race(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(t *testing.T) {
			delegate := &aiUpDownCounter{}
			f := func(int64) { _ = delegate.Unwrap() }
			testInt64Race(f, delegate.setDelegate)
		})

		t.Run("Gauge", func(t *testing.T) {
			delegate := &aiGauge{}
			f := func(int64) { _ = delegate.Unwrap() }
			testInt64Race(f, delegate.setDelegate)
		})
	})
}

func TestSyncInstrumentSetDelegateRace(t *testing.T) {
	// Float64 Instruments
	t.Run("Float64", func(t *testing.T) {
		t.Run("Counter", func(t *testing.T) {
			delegate := &sfCounter{}
			f := func(v float64) { delegate.Add(context.Background(), v) }
			testFloat64Race(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(t *testing.T) {
			delegate := &sfUpDownCounter{}
			f := func(v float64) { delegate.Add(context.Background(), v) }
			testFloat64Race(f, delegate.setDelegate)
		})

		t.Run("Histogram", func(t *testing.T) {
			delegate := &sfHistogram{}
			f := func(v float64) { delegate.Record(context.Background(), v) }
			testFloat64Race(f, delegate.setDelegate)
		})
	})

	// Int64 Instruments

	t.Run("Int64", func(t *testing.T) {
		t.Run("Counter", func(t *testing.T) {
			delegate := &siCounter{}
			f := func(v int64) { delegate.Add(context.Background(), v) }
			testInt64Race(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(t *testing.T) {
			delegate := &siUpDownCounter{}
			f := func(v int64) { delegate.Add(context.Background(), v) }
			testInt64Race(f, delegate.setDelegate)
		})

		t.Run("Histogram", func(t *testing.T) {
			delegate := &siHistogram{}
			f := func(v int64) { delegate.Record(context.Background(), v) }
			testInt64Race(f, delegate.setDelegate)
		})
	})
}

type testCountingFloatInstrument struct {
	count int

	instrument.Float64Observable
	embedded.Float64Counter
	embedded.Float64UpDownCounter
	embedded.Float64Histogram
	embedded.Float64ObservableCounter
	embedded.Float64ObservableUpDownCounter
	embedded.Float64ObservableGauge
}

func (i *testCountingFloatInstrument) observe() {
	i.count++
}
func (i *testCountingFloatInstrument) Add(context.Context, float64, ...instrument.AddOption) {
	i.count++
}
func (i *testCountingFloatInstrument) Record(context.Context, float64, ...instrument.RecordOption) {
	i.count++
}

type testCountingIntInstrument struct {
	count int

	instrument.Int64Observable
	embedded.Int64Counter
	embedded.Int64UpDownCounter
	embedded.Int64Histogram
	embedded.Int64ObservableCounter
	embedded.Int64ObservableUpDownCounter
	embedded.Int64ObservableGauge
}

func (i *testCountingIntInstrument) observe() {
	i.count++
}
func (i *testCountingIntInstrument) Add(context.Context, int64, ...instrument.AddOption) {
	i.count++
}
func (i *testCountingIntInstrument) Record(context.Context, int64, ...instrument.RecordOption) {
	i.count++
}
