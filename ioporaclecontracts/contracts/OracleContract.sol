pragma solidity >=0.4.22 <0.8.0;

contract OracleContract {
    struct IopNode {
        address addr;
        string ipAddr;
        uint256 index;
    }

    mapping(address => IopNode) private iopNodes;
    address[] private iopNodeIndices;

    bool public result;
    event VerifyTransactionLog(uint256 id, string hash, uint256 confirmations);

    function registerIopNode(string calldata _ipAddr) external payable {
        require(!iopNodeIsRegistered(msg.sender), "already registered");
        IopNode storage iopNode = iopNodes[msg.sender];
        iopNode.addr = msg.sender;
        iopNode.ipAddr = _ipAddr;
        iopNode.index = iopNodeIndices.length;
        iopNodeIndices.push(iopNode.addr);
    }

    function iopNodeIsRegistered(address _addr) public view returns (bool) {
        if (iopNodeIndices.length == 0) return false;
        return (iopNodeIndices[iopNodes[_addr].index] == _addr);
    }

    function findIopNodeByAddress(address _addr)
        external
        view
        returns (address, string memory)
    {
        require(iopNodeIsRegistered(_addr), "not found");
        IopNode memory iopNode = iopNodes[_addr];
        return (iopNode.addr, iopNode.ipAddr);
    }

    function findIopNodeByIndex(uint256 _index)
        external
        view
        returns (address, string memory)
    {
        require(_index >= 0 && _index < iopNodeIndices.length, "not found");
        IopNode memory iopNode = iopNodes[iopNodeIndices[_index]];
        return (iopNode.addr, iopNode.ipAddr);
    }

    function countIopNodes() external view returns (uint256) {
        return iopNodeIndices.length;
    }

    function verifyTransaction(string calldata _hash, uint256 _confirmations)
        external
    {
        emit VerifyTransactionLog(1, _hash, _confirmations);
    }

    function verifyTransactionResult(bool _result) external {
        result = _result;
    }
}
