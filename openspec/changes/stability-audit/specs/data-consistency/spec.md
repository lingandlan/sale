## ADDED Requirements

### Requirement: C recharge compensation flow
The `CreateCRecharge` service method SHALL implement a compensation-based flow:
1. Create CRecharge record with `status = "pending"`
2. Deduct center balance (database operation)
3. Call WSY AddIntegral API
4. If AddIntegral succeeds → update CRecharge status to "success" and set balance_after
5. If AddIntegral fails → keep CRecharge status as "pending", log the error, and return a business error

The CRecharge model SHALL support a `status` field with values: "pending", "success", "failed".

#### Scenario: Successful C recharge
- **WHEN** C recharge is submitted with valid parameters and WSY AddIntegral succeeds
- **THEN** the system SHALL:
  1. Create CRecharge record with status "success"
  2. Deduct center balance
  3. Add points to member via WSY API
  4. Update balance_after with the actual member balance after points added
  5. Return the recharge record

#### Scenario: WSY AddIntegral fails
- **WHEN** C recharge is submitted and WSY AddIntegral returns an error
- **THEN** the system SHALL:
  1. Keep CRecharge record with status "pending"
  2. Center balance already deducted (not rolled back)
  3. Log the error with full context (memberPhone, points, centerID)
  4. Return a business error to the caller

#### Scenario: CRecharge record creation fails
- **WHEN** C recharge is submitted and the database insert for CRecharge fails
- **THEN** the system SHALL return the database error without deducting center balance or calling WSY API

### Requirement: Card status transition transaction
The `transitionCardStatus` operation SHALL execute the card status update AND the card transaction record creation within a single GORM transaction. The repository SHALL provide a `TransitionCardStatusTX` method that performs both operations atomically.

#### Scenario: Successful card status transition
- **WHEN** a card status is transitioned from "in_stock" to "frozen"
- **THEN** the system SHALL:
  1. Update card status in the database
  2. Create a card_transaction audit record
  3. Both operations SHALL succeed or both SHALL fail (rollback)

#### Scenario: Transaction record creation fails
- **WHEN** card status update succeeds but the audit record insert fails
- **THEN** the entire transaction SHALL roll back, and the card status SHALL remain unchanged

### Requirement: Batch card creation transaction wrapping
After `BatchCreateCards` succeeds, the subsequent loop creating card transaction records SHALL be executed within a transaction. If any transaction record creation fails, all previously created records in that batch SHALL roll back.

#### Scenario: Batch of 50 cards, transaction record for card 30 fails
- **WHEN** BatchCreateCards inserts 50 cards successfully, and the 30th card_transaction insert fails
- **THEN** all 30 card_transaction records SHALL be rolled back, but the 50 card records SHALL remain (they were inserted in a separate successful transaction)
- **THEN** the service SHALL return the error from the failed transaction record insert
