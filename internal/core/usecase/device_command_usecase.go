package usecase

import (
	"project-home-iot/internal/core/domain"

)

type CommandUsecase interface {
	SendCommand(cmd *domain.DeviceCommand) error
}

type commandUsecase struct {
    deviceRepo       domain.DeviceRepository
    widgetRepository domain.WidgetRepository
    commander        domain.DeviceCommander
    recorder      domain.Recorder
}

func NewCommandUsecase(dr domain.DeviceRepository, wr domain.WidgetRepository, dc domain.DeviceCommander,rd domain.Recorder) CommandUsecase {
    return &commandUsecase{
        deviceRepo:       dr,
        widgetRepository: wr,
        commander:        dc,
        recorder:      rd,
    }
}

func (u *commandUsecase) SendCommand(cmd *domain.DeviceCommand) error {
    device, err := u.deviceRepo.FindByWidgetID(cmd.WidgetID)
    if err != nil {
        return err
    }

    resp, err := u.commander.RequestCommand(device.Topic, cmd)
    if err != nil {
        return err
    }

    u.recorder.RecordLog(cmd.WidgetID, "command", resp.Value)

    return u.widgetRepository.UpdateValue(cmd.WidgetID, resp.Value)
}

// func (u *commandUsecase) SendCommand(cmd *domain.DeviceCommand) error {
//     device, err := u.deviceRepo.FindByWidgetID(cmd.WidgetID)
//     if err != nil {
//         return err
//     }

//     // 1. สร้าง Topic เฉพาะสำหรับรอรับ Reply ของ Request นี้เท่านั้น
//     // สมมติว่า cmd.CorrelationID มีค่า เช่น "req-123"
//     replyTopic := fmt.Sprintf("devices/reply/%s", cmd.CorrelationID)
    
// 	fmt.Printf("🔄 Sending command to device on topic: %s\n", device.Topic)
//     // 2. เตรียม Channel มารอรับข้อมูล (กำหนดขนาด 1 เพื่อไม่ให้บล็อกตอนส่งเข้า)
//     ch := make(chan dto.CommandResponse, 1)

//     u.mqttClient.Subscribe(replyTopic, 1, func(c mqtt.Client, m mqtt.Message) {
//         var resp dto.CommandResponse
//         if err := json.Unmarshal(m.Payload(), &resp); err != nil {
//             fmt.Printf("Unmarshal error: %v\n", err)
//             return
//         }

//         fmt.Printf("📥 Received matching reply on %s\n", replyTopic)
//         ch <- resp
        
//         c.Unsubscribe(replyTopic)
//     })

//     // 4. ส่งคำสั่งออกไป
//     payload, _ := json.Marshal(cmd)
//     token := u.mqttClient.Publish(device.Topic, 1, false, payload)
//     token.Wait()
//     if token.Error() != nil {
//         return token.Error()
//     }
//     fmt.Println("📤 MQTT published to:", device.Topic)

//     // 5. รอรับ Reply พร้อมระบบ Timeout
//     select {
//     case receivedResp := <-ch:
// 		u.widgetRepository.UpdateValue(cmd.WidgetID, receivedResp.Value)
//         fmt.Printf("✅ Success: %+v\n", receivedResp)
//         // h.widgetUsecase.UpdateValue(receivedResp.WidgetID, receivedResp.Value)
//         return nil

//     case <-time.After(30 * time.Second): // รอ 10 วินาที ถ้าไม่ตอบให้คืน Error
//         u.mqttClient.Unsubscribe(replyTopic) // อย่าลืมเลิกฟังถ้า timeout
//         return fmt.Errorf("timeout: device did not respond on %s within 10s", replyTopic)
//     }
// }

// func (u *commandUsecase) SendCommand(cmd *domain.DeviceCommand) error {

// 	device, err := u.deviceRepo.FindByWidgetID(cmd.WidgetID)
// 	if err != nil {
// 		return err
// 	}
// 	payload, _ := json.Marshal(cmd)

// 	targetID := cmd.CorrelationID

// 	token := u.mqttClient.Publish(
// 		device.Topic,
// 		1,
// 		false,
// 		payload,
// 	)

// 	fmt.Println("📤 MQTT publish")
// 	fmt.Println("Topic:", device.Topic)
// 	fmt.Println("Payload:", string(payload))

// 	token.Wait()

// 	ch := make(chan dto.CommandResponse)

// 	u.mqttClient.Subscribe("devices/reply/+", 1, func(c mqtt.Client, m mqtt.Message) {

// 		var resp dto.CommandResponse
// 		json.Unmarshal(m.Payload(), &resp)

		

// 		fmt.Printf("Received reply for correlation ID %s: Widget %d new value %d\n", resp.CorrelationID, resp.WidgetID, resp.Value)
// 		// h.widgetUsecase.UpdateValue(resp.WidgetID, resp.Value)
// 		if resp.CorrelationID == targetID {
//             fmt.Printf("✅ Matched ID: %s\n", resp.CorrelationID)
//             ch <- resp // ส่งเข้า Channel เฉพาะตัวที่ ID ตรงกัน
//         } else {
//             fmt.Printf("❌ Ignored ID: %s (Waiting for %s)\n", resp.CorrelationID, targetID)
//         }
// 		fmt.Println("📥 MQTT waiting....")
// 		ch <- resp
// 		fmt.Println("📥 MQTT continue")
// 	})

// 	receivedResp := <-ch
// 	fmt.Printf("Final received response: %+v\n", receivedResp)


// 	return token.Error()
// }
