## ADDED Requirements

### Requirement: UserInfo persisted to localStorage
User store SHALL persist `userInfo` to localStorage key `user_info` whenever it changes, and restore it on store initialization.

#### Scenario: Page refresh restores userInfo
- **WHEN** user is logged in (token + userInfo in store), page is refreshed
- **THEN** userInfo is restored from localStorage, no additional API call needed

#### Scenario: Logout clears persisted userInfo
- **WHEN** user logs out
- **THEN** localStorage `user_info` is removed along with `access_token` and `refresh_token`

#### Scenario: Corrupted localStorage data
- **WHEN** localStorage `user_info` contains invalid JSON
- **THEN** store ignores it gracefully, userInfo remains null
