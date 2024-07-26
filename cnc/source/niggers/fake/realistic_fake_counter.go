package fake

import (
	"math/rand"
	"time"
)

func realisticFakeCounter(deviceType string, deviceArch string, minAmount int, maxAmount int) {
	counter := rand.Intn(maxAmount-minAmount+1) + minAmount
	isConnected := true

	for {
		sleepDuration := rand.Intn(10) + 1
		time.Sleep(time.Duration(sleepDuration) * time.Second)
		if isConnected {
			if counter >= minAmount && counter <= maxAmount {
				if rand.Float32() < 0.5 {
					counter += rand.Intn(6) + 5
				} else {
					counter -= rand.Intn(6) + 5
				}
			} else {
				counter = rand.Intn(maxAmount-minAmount+1) + minAmount
			}

			fakeBots[deviceType] = Device{
				Name:    deviceType,
				Arch:    deviceArch,
				Current: counter,
				Minimum: minAmount,
				Maximum: maxAmount,
			}
		} else {
			sleepDuration := rand.Intn(10) + 1
			time.Sleep(time.Duration(sleepDuration) * time.Second)
		}

		if rand.Float32() < 0.05 {
			isConnected = !isConnected
		}
	}
}
