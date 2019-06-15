UPDATE user SET
  block_remaining = IF(CAST(block_remaining as SIGNED) - ? < 0, 0, block_remaining - ?)
WHERE user = ?