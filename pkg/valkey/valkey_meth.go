package valkey

import (
	"context"
	"encoding/json"
	"fmt"
	"fyc/config"

	"github.com/rs/zerolog/log"
	"github.com/valkey-io/valkey-go"
)

type ValkeyStrct struct {
	client valkey.Client
}

func (v *ValkeyStrct) Valkey_Connect() {
	var LinkConnect = fmt.Sprintf("%s:%d", config.Configvar.Valkey.Host, config.Configvar.Valkey.Port)

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{LinkConnect},
	})
	if err != nil {
		log.Err(err).Msg("Failed to create Valkey client")
		return
	}

	v.client = client
	log.Info().Str("Host", LinkConnect).Msg("- - - - - - - CONNECTED TO VALKEY SERVER - - - - - - - -")
}

func (v *ValkeyStrct) Valkey_Close() {
	if v.client != nil {
		v.client.Close()
		log.Info().Msg("Disconnected from Valkey server")
	}
}
func (v *ValkeyStrct) PublishMessage(ctx context.Context, val interface{}) {
	channel := config.Configvar.Valkey.Channel

	data, err := json.Marshal(val)
	if err != nil {
		log.Err(err).Msg("Error marshaling data to JSON:")
		return
	}

	err = v.client.Do(ctx, v.client.B().Publish().Channel(channel).Message(string(data)).Build()).Error()
	if err != nil {
		log.Err(err).Msgf("Error publishing to channel : %s  / Erorr :", channel)
	}

	log.Info().Str("channel", channel).Msg("Successfully published message")
}

func (v *ValkeyStrct) Valkey_Setter_Data(ctx context.Context, key string, value string) bool {
	log.Debug().Msg("------------- Setting Data into Valkey DB ----------")

	resp, err := v.client.Do(ctx, v.client.B().Set().Key(key).Value(value).Build()).ToString()
	if err != nil {
		log.Error().Err(err).Msgf("Error setting value %s to Key %s", value, key)
		return false
	}

	if resp != "OK" {
		log.Error().Str("Response", resp).Msgf("Unexpected response setting value %s to Key %s", value, key)
		return false
	}

	log.Info().Str("Value", value).Str("Key", key).Msg("Value set successfully in Valkey DB")
	return true
}

func (v *ValkeyStrct) Valkey_Getter_Data(ctx context.Context, key string) (int64, bool) {

	log.Debug().Msg("------------- Retrieving Data from Valkey DB ----------")

	val, err := v.client.Do(ctx, v.client.B().Get().Key(key).Build()).AsInt64()
	if err != nil {
		log.Error().Err(err).Msgf("Error retrieving value for Key %s", key)
		return 0, false
	}

	log.Info().Int64("Value", val).Str("Key", key).Msg("Value retrieved successfully from Valkey DB")
	return val, true
}

func (v *ValkeyStrct) Valkey_Incr_Data(ctx context.Context, key string) bool {
	log.Debug().Msg("------------- Increment Data in Valkey DB ----------")

	if v.client == nil {
		log.Error().Msg("Valkey client is not initialized. Call Valkey_Connect first.")
		return false
	}

	val, err := v.client.Do(ctx, v.client.B().Incr().Key(key).Build()).AsInt64()
	if err != nil {
		log.Error().Err(err).Msgf("Error incrementing value for Key %s", key)
		return false
	}

	log.Info().Int64("New Value", val).Str("Key", key).Msg("New value incremented in Valkey DB")
	return true
}

func (v *ValkeyStrct) Valkey_Decr_Data(ctx context.Context, key string) bool {

	log.Debug().Msg("------------- Decremant Data into Valkey DB ----------")
	val, err := v.client.Do(ctx, v.client.B().Decr().Key(key).Build()).AsInt64()
	if err != nil {
		log.Error().Err(err).Msgf("Error Decre new value to Key %s", key)
		return false
	}

	log.Info().Int64("New Value", val).Str("Key", key).Msg("New value Decremented into Valkey DB")
	return true
}
