//init test queries
export const GET_COUNTRIES = `query GetCountries {
  getCountries {
    countries {
      name
      code
    }
  }
}`;

export const GET_TIMEZONES = `query GetTimezones {
  getTimezones {
    timezones {
      value
      abbr
      offset
      isdst
      text
      utc
    }
  }
}`;

//invitation test queries
export const CREATE_INVITATION = `mutation CreateInvitation($data: CreateInvitationInputDto!) {
    createInvitation(data: $data) {
      email
      expireAt
      id
      name
      role
      link
      userId
      status
    }
  }`;

export const GET_INVITATION = `query GetInvitation($getInvitationId: String!) {
  getInvitation(id: $getInvitationId) {
    email
    expireAt
    id
    name
    role
    link
    userId
    status
  }
}`;

export const GET_ORG_INVITATION = `query GetInvitationsByOrg {
  getInvitationsByOrg {
    invitations {
      email
      expireAt
      id
      name
      role
      link
      userId
      status
    }
  }
}`;

export const GET_EMAIL_INVITATIONS = `query GetInvitations($email: String!) {
  getInvitations(email: $email) {
    email
    expireAt
    id
    name
    role
    link
    userId
    status
  }
}`;

export const UPDATE_INVITATION = `mutation UpdateInvitation($data: UpateInvitationInputDto!) {
  updateInvitation(data: $data) {
    id
  }
}`;

export const DELETE_INVITATION = `mutation DeleteInvitation($deleteInvitationId: String!) {
    deleteInvitation(id: $deleteInvitationId) {
      id
    }
  }`;

//member test queries
export const GET_MEMBERS = `query GetMembers {
  getMembers {
    members {
      userId
      name
      email
      memberId
      isDeactivated
      role
      memberSince
    }
  }
}`;

export const ADD_MEMBER = `mutation AddMember($data: AddMemberInputDto!) {
  addMember(data: $data) {
    userId
    memberId
    isDeactivated
    memberSince
    role
  }
}`;

export const GET_MEMBER = `query GetMemberByUserId($userId: String!) {
  getMemberByUserId(userId: $userId) {
    userId
    name
    email
    memberId
    isDeactivated
    role
    memberSince
  }
}`;

export const GET_MEMBER_BY_ID = `query GetMember($getMemberId: String!) {
    getMember(id: $getMemberId) {
      userId
      memberId
      isDeactivated
      role
      memberSince
    }
  }`;

export const UPDATE_MEMBER = `mutation UpdateMember($data: UpdateMemberInputDto!, $memberId: String!) {
  updateMember(data: $data, memberId: $memberId) {
    success
  }
}`;

export const REMOVE_MEMBER = `mutation RemoveMember($removeMemberId: String!) {
  removeMember(id: $removeMemberId) {
    success
  }
}`;

//network test queries
export const ADD_NETWORK = `mutation AddNetwork($data: AddNetworkInputDto!) {
  addNetwork(data: $data) {
    id
    name
    isDefault
    budget
    overdraft
    trafficPolicy
    isDeactivated
    paymentLinks
    createdAt
    countries
    networks
  }
}`;

export const GET_NETWORKS = `query GetNetworks {
  getNetworks {
    networks {
      id
      name
      isDefault
      budget
      overdraft
      trafficPolicy
      isDeactivated
      paymentLinks
      createdAt
      countries
      networks
    }
  }
}`;

export const GET_NETWORK = `query GetNetwork($networkId: String!) {
  getNetwork(networkId: $networkId) {
    id
    name
    isDefault
    budget
    overdraft
    trafficPolicy
    isDeactivated
    paymentLinks
    createdAt
    countries
    networks
  }
}`;

export const SET_DEFAULT = `mutation SetDefaultNetwork($data: SetDefaultNetworkInputDto!) {
  setDefaultNetwork(data: $data) {
    success
  }
}`;

// node test queries
export const ADD_NODE = `mutation AddNode($data: AddNodeInput!) {
  addNode(data: $data) {
    id
    name
    orgId
    type
    attached {
      id
      name
      orgId
      type
    }
    site {
      nodeId
      siteId
      networkId
      addedAt
    }
    status {
      connectivity
      state
    }
  }
}`;

export const ADD_NODE_TO_SITE = `mutation AddNodeToSite($data: AddNodeToSiteInput!) {
  addNodeToSite(data: $data) {
    success
  }
}`;

export const ATTACH_NODE = `mutation AttachNode($data: AttachNodeInput!) {
  attachNode(data: $data) {
    success
  }
}`;

export const GET_APPS_CHANGE = `query GetAppsChangeLog($data: NodeAppsChangeLogInput!) {
  getAppsChangeLog(data: $data) {
    logs {
      version
      date
    }
    type
  }
}`;

export const GET_NODE = `query GetNode($data: NodeInput!) {
  getNode(data: $data) {
    id
    name
    orgId
    type
    attached {
      id
      name
      orgId
      type
    }
    site {
      nodeId
      siteId
      networkId
      addedAt
    }
    status {
      connectivity
      state
    }
  }
}`;

export const GET_NODES = `query GetNodes($data: GetNodesInput!) {
  getNodes(data: $data) {
    nodes {
      id
      name
      orgId
      type
      attached {
        id
        name
        orgId
        type
      }
      site {
        nodeId
        siteId
        networkId
        addedAt
      }
      status {
        connectivity
        state
      }
    }
  }
}`;

export const GET_NODE_APPS = `query GetNodeApps($data: NodeAppsChangeLogInput!) {
  getNodeApps(data: $data) {
    apps {
      name
      date
      version
      cpu
      memory
      notes
    }
    type
  }
}`;

export const GET_NODE_LOCATION = `query GetNodeLocation($data: NodeInput!) {
  getNodeLocation(data: $data) {
    id
    lat
    lng
    state
  }
}`;

export const GET_NETWORK_NODES = `query GetNodesByNetwork($networkId: String!) {
  getNodesByNetwork(networkId: $networkId) {
    nodes {
      id
      name
      orgId
      type
      attached {
        id
        name
        orgId
        type
      }
      site {
        nodeId
        siteId
        networkId
        addedAt
      }
      status {
        connectivity
        state
      }
    }
  }
}`;

export const GET_NODES_LOCATION = `query GetNodesLocation($data: NodesInput!) {
  getNodesLocation(data: $data) {
    networkId
    nodes {
      id
      lat
      lng
      state
    }
  }
}`;

export const UPDATE_NODE = `mutation UpdateNode($data: UpdateNodeInput!) {
  updateNode(data: $data) {
    id
    name
    orgId
    type
    attached {
      id
      name
      orgId
      type
    }
    site {
      nodeId
      siteId
      networkId
      addedAt
    }
    status {
      connectivity
      state
    }
  }
}`;

export const UPDATE_NODE_STATE = `mutation UpdateNodeState($data: UpdateNodeStateInput!) {
  updateNodeState(data: $data) {
    id
    name
    orgId
    type
    attached {
      id
      name
      orgId
      type
    }
    site {
      nodeId
      siteId
      networkId
      addedAt
    }
    status {
      connectivity
      state
    }
  }
}`;

export const DETACH_NODE = `mutation DetachhNode($data: NodeInput!) {
  detachhNode(data: $data) {
    success
  }
}`;

export const DELETE_NODE = `mutation DeleteNodeFromOrg($data: NodeInput!) {
  deleteNodeFromOrg(data: $data) {
    id
  }
}`;

// org test queries
export const GET_ORG = `query GetOrg {
  getOrg {
    id
    name
    owner
    certificate
    isDeactivated
    createdAt
  }
}`;

export const GET_ORGS = `query GetOrgs {
  getOrgs {
    user
    ownerOf {
      id
      name
      owner
      certificate
      isDeactivated
      createdAt
    }
    memberOf {
      id
      name
      owner
      certificate
      isDeactivated
      createdAt
    }
  }
}`;

// package test queries
export const ADD_PACKAGE = `mutation AddPackage($data: AddPackageInputDto!) {
  addPackage(data: $data) {
    uuid
    name
    active
    duration
    simType
    createdAt
    deletedAt
    updatedAt
    smsVolume
    dataVolume
    voiceVolume
    ulbr
    dlbr
    type
    dataUnit
    voiceUnit
    messageUnit
    flatrate
    currency
    from
    to
    country
    provider
    apn
    ownerId
    amount
    rate {
      sms_mo
      sms_mt
      data
      amount
    }
    markup {
      baserate
      markup
    }
  }
}`;

export const GET_PACKAGE = `query GetPackage($packageId: String!) {
  getPackage(packageId: $packageId) {
    uuid
    name
    active
    duration
    simType
    createdAt
    deletedAt
    updatedAt
    smsVolume
    dataVolume
    voiceVolume
    ulbr
    dlbr
    type
    dataUnit
    voiceUnit
    messageUnit
    flatrate
    currency
    from
    to
    country
    provider
    apn
    ownerId
    amount
    rate {
      sms_mo
      sms_mt
      data
      amount
    }
    markup {
      baserate
      markup
    }
  }
}`;

export const GET_PACKAGES = `query GetPackages {
  getPackages {
    packages {
      uuid
      name
      active
      duration
      simType
      createdAt
      deletedAt
      updatedAt
      smsVolume
      dataVolume
      voiceVolume
      ulbr
      dlbr
      type
      dataUnit
      voiceUnit
      messageUnit
      flatrate
      currency
      from
      to
      country
      provider
      apn
      ownerId
      amount
      rate {
        sms_mo
        sms_mt
        data
        amount
      }
      markup {
        baserate
        markup
      }
    }
  }
}`;

export const DELETE_PACKAGE = `mutation DeletePackage($packageId: String!) {
  deletePackage(packageId: $packageId) {
    uuid
  }
}`;

// sim test queries
export const UPLOAD_SIMS = `mutation UploadSims($data: UploadSimsInputDto!) {
  uploadSims(data: $data) {
    iccid
  }
}`;

export const ADD_SIM_PACKAGE = `mutation AddPackageToSim($data: AddPackageToSimInputDto!) {
  addPackageToSim(data: $data) {
    packageId
  }
}`;

export const ALLOCATE_SIM = `mutation AllocateSim($data: AllocateSimInputDto!) {
  allocateSim(data: $data) {
    id
    subscriber_id
    network_id
    package {
      id
      packageId
      startDate
      endDate
      isActive
    }
    iccid
    msisdn
    imsi
    type
    status
    is_physical
    traffic_policy
    firstActivatedOn
    lastActivatedOn
    activationsCount
    deactivationsCount
    allocated_at
    sync_status
  }
}`;

export const GET_SIM = `query GetSim($data: GetSimInputDto!) {
  getSim(data: $data) {
    activationCode
    createdAt
    iccid
    id
    isAllocated
    isPhysical
    msisdn
    qrCode
    simType
    smapAddress
  }
}`;

export const GET_SIM_STATS = `query GetSimPoolStats($type: String!) {
  getSimPoolStats(type: $type) {
    total
    available
    consumed
    failed
    esim
    physical
  }
}`;

export const GET_DATA_USAGE = `query GetDataUsage($simId: String!) {
  getDataUsage(simId: $simId) {
    usage
  }
}`;

export const GET_SIM_PACKAGES = `query GetPackagesForSim($data: GetPackagesForSimInputDto!) {
  getPackagesForSim(data: $data) {
    sim_id
    packages {
      id
      package_id
      start_date
      end_date
      is_active
    }
  }
}`;

export const GET_SUBSCRIBER_SIM = `query GetSimsBySubscriber($data: GetSimBySubscriberInputDto!) {
  getSimsBySubscriber(data: $data) {
    subscriber_id
    sims {
      activationCode
      createdAt
      iccid
      id
      isAllocated
      isPhysical
      msisdn
      qrCode
      simType
      smapAddress
    }
  }
}`;

export const DELETE_SIM = `mutation DeleteSim($data: DeleteSimInputDto!) {
  deleteSim(data: $data) {
    simId
  }
}`;

// subscriber test queries
export const GET_SUBSCRIBER = `query GetSubscriber($subscriberId: String!) {
    getSubscriber(subscriberId: $subscriberId) {
      uuid
      address
      dob
      email
      firstName
      lastName
      gender
      idSerial
      networkId
      phone
      proofOfIdentification
      sim {
        id
        subscriberId
        networkId
        iccid
        msisdn
        imsi
        type
        status
        firstActivatedOn
        lastActivatedOn
        activationsCount
        deactivationsCount
        allocatedAt
        isPhysical
        package
      }
    }
  }`;

export const GET_SUBSCRIBER_METRICS = `query GetSubscriberMetricsByNetwork($networkId: String!) {
  getSubscriberMetricsByNetwork(networkId: $networkId) {
    total
    active
    inactive
    terminated
  }
}`;

export const GET_SUBSCRIBERS = `query GetSubscribersByNetwork($networkId: String!) {
  getSubscribersByNetwork(networkId: $networkId) {
    subscribers {
      uuid
      address
      dob
      email
      firstName
      lastName
      gender
      idSerial
      networkId
      phone
      proofOfIdentification
      sim {
        id
        subscriberId
        networkId
        iccid
        msisdn
        imsi
        type
        status
        firstActivatedOn
        lastActivatedOn
        activationsCount
        deactivationsCount
        allocatedAt
        isPhysical
        package
      }
    }
  }
}`;

export const UPDATE_SUBSCRIBER = `mutation UpdateSubscriber($data: UpdateSubscriberInputDto!, $subscriberId: String!) {
  updateSubscriber(data: $data, subscriberId: $subscriberId) {
    success
  }
}`;

export const DELETE_SUBSCRIBER = `mutation DeleteSubscriber($subscriberId: String!) {
  deleteSubscriber(subscriberId: $subscriberId) {
    success
  }
}`;

// user test queries
export const GET_USER = `query GetUser($userId: String!) {
  getUser(userId: $userId) {
    name
    email
    uuid
    phone
    isDeactivated
    authId
    registeredSince
  }
}`;

export const WHO_AM_I = `query Whoami {
  whoami {
    user {
      name
      email
      uuid
      phone
      isDeactivated
      authId
      registeredSince
    }
    ownerOf {
      id
      name
      owner
      certificate
      isDeactivated
      createdAt
    }
    memberOf {
      id
      name
      owner
      certificate
      isDeactivated
      createdAt
    }
  }
}`;
