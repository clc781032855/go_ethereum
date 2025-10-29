// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
    uint256 private count;
    address public owner;

    // 事件声明
    event Incremented(uint256 newCount);
    event Decremented(uint256 newCount);
    event Reset(uint256 newCount);

    // 构造函数
    constructor() {
        owner = msg.sender;
        count = 0;
    }

    // 获取当前计数
    function getCount() public view returns (uint256) {
        return count;
    }

    // 增加计数
    function increment() public {
        count += 1;
        emit Incremented(count);
    }

    // 减少计数
    function decrement() public {
        require(count > 0, "Count cannot be negative");
        count -= 1;
        emit Decremented(count);
    }

    // 重置计数
    function reset() public {
        count = 0;
        emit Reset(count);
    }
}
