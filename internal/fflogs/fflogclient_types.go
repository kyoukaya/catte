package fflogs

type RawBuffEvent struct {
	Timestamp           int64    `json:"timestamp"`
	Type                BuffType `json:"type"`
	SourceID            *int64   `json:"sourceID,omitempty"`
	TargetID            *int64   `json:"targetID,omitempty"`
	AbilityGameID       *int64   `json:"abilityGameID,omitempty"`
	HitType             *int64   `json:"hitType,omitempty"`
	Amount              *int64   `json:"amount,omitempty"`
	Multiplier          *float64 `json:"multiplier,omitempty"`
	PacketID            *int64   `json:"packetID,omitempty"`
	Stack               *int64   `json:"stack,omitempty"`
	DirectHit           *bool    `json:"directHit,omitempty"`
	Value               *int64   `json:"value,omitempty"`
	Bars                *int64   `json:"bars,omitempty"`
	Melee               *bool    `json:"melee,omitempty"`
	Absorbed            *int64   `json:"absorbed,omitempty"`
	Absorb              *int64   `json:"absorb,omitempty"`
	Overheal            *int64   `json:"overheal,omitempty"`
	Tick                *bool    `json:"tick,omitempty"`
	FinalizedAmount     *float64 `json:"finalizedAmount,omitempty"`
	Simulated           *bool    `json:"simulated,omitempty"`
	ExpectedAmount      *int64   `json:"expectedAmount,omitempty"`
	ExpectedCritRate    *int64   `json:"expectedCritRate,omitempty"`
	ActorPotencyRatio   *float64 `json:"actorPotencyRatio,omitempty"`
	GuessAmount         *float64 `json:"guessAmount,omitempty"`
	DirectHitPercentage *float64 `json:"directHitPercentage,omitempty"`
	TargetInstance      *int64   `json:"targetInstance,omitempty"`
}

type BuffType string

const (
	Applybuff        BuffType = "applybuff"
	Applybuffstack   BuffType = "applybuffstack"
	Applydebuff      BuffType = "applydebuff"
	Applydebuffstack BuffType = "applydebuffstack"
	Begincast        BuffType = "begincast"
	Calculateddamage BuffType = "calculateddamage"
	Calculatedheal   BuffType = "calculatedheal"
	Cast             BuffType = "cast"
	Damage           BuffType = "damage"
	Heal             BuffType = "heal"
	Limitbreakupdate BuffType = "limitbreakupdate"
	Refreshbuff      BuffType = "refreshbuff"
	Refreshdebuff    BuffType = "refreshdebuff"
	Removebuff       BuffType = "removebuff"
	Removebuffstack  BuffType = "removebuffstack"
	Removedebuff     BuffType = "removedebuff"
)

type EventDataType string

const (
	EDTypeAll           EventDataType = "All"
	EDTypeBuffs         EventDataType = "Buffs"
	EDTypeCasts         EventDataType = "Casts"
	EDTypeCombatantInfo EventDataType = "CombatantInfo"
	EDTypeDamageDone    EventDataType = "DamageDone"
	EDTypeDamageTaken   EventDataType = "DamageTaken"
	EDTypeDeaths        EventDataType = "Deaths"
	EDTypeDebuffs       EventDataType = "Debuffs"
	EDTypeDispels       EventDataType = "Dispels"
	EDTypeHealing       EventDataType = "Healing"
	EDTypeInterrupts    EventDataType = "Interrupts"
	EDTypeResources     EventDataType = "Resources"
	EDTypeSummons       EventDataType = "Summons"
	EDTypeThreat        EventDataType = "Threat"
)

type HostilityType string

const (
	Friendlies HostilityType = "Friendlies"
	Enemies    HostilityType = "Enemies"
)
