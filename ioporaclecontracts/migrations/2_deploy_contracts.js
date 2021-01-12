const ECDSAOracleContract = artifacts.require("ECDSAOracleContract");
const BLSOracleContract = artifacts.require("BLSOracleContract");
const MerkleOracleContract = artifacts.require("MerkleOracleContract");
const RegistryContract = artifacts.require("RegistryContract");

module.exports = function(deployer) {
  deployer.deploy(RegistryContract).then(function() {
    return deployer
      .deploy(ECDSAOracleContract, RegistryContract.address)
      .then(function() {
        return deployer
          .deploy(MerkleOracleContract, RegistryContract.address)
          .then(function() {
            return deployer.deploy(BLSOracleContract, RegistryContract.address);
          });
      });
  });
};
