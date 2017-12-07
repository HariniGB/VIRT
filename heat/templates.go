package heat

const (
	Network = `heat_template_version: '2013-05-23'
parameters:
  net_name:
    type: string
    description: Name of network
  subnet_name:
    type: string
    description: Name of subnet
  router_name:
    type: string
    description: Name of router
  cidr:
    type: string
    description: CIDR for subnet

resources:
  server_security_group:
    type: OS::Neutron::SecurityGroup
    properties:
      description: Add security group rules for server
      name: ssh-ping
      rules:
        - remote_ip_prefix: 0.0.0.0/0
          protocol: tcp
          port_range_min: 22
          port_range_max: 22
        - remote_ip_prefix: 0.0.0.0/0
          protocol: icmp

  private_net:
    properties: {name: { get_param: net_name }}
    type: OS::Neutron::Net
  private_subnet:
    properties:
      name: { get_param: subnet_name }
      cidr: { get_param: cidr }
      network_id: {get_resource: private_net}
    type: OS::Neutron::Subnet
  router1:
    properties:
      name: { get_param: router_name }
      external_gateway_info: {network: provider}
    type: OS::Neutron::Router
  router1_interface:
    properties:
      router_id: {get_resource: router1}
      subnet_id: {get_resource: private_subnet}
    type: OS::Neutron::RouterInterface`

    Cirros = `heat_template_version: '2013-05-23'

parameters:
  name:
    type: string
    description: Name of the VM

  net_name:
    type: string
    description: Name of network

  subnet_name:
    type: string
    description: Name of subnet

  flavor_name:
    type: string
    description: Name of the flavor

resources:
  server1_port:
    type: OS::Neutron::Port
    properties:
      network_id: { get_param: net_name }
      security_groups: [ ssh-ping, default ]
      fixed_ips:
        - subnet_id: { get_param: subnet_name }

  server1_floating_ip:
    type: OS::Neutron::FloatingIP
    properties:
      floating_network_id: provider
      port_id: { get_resource: server1_port }

  server1:
    type: OS::Nova::Server
    properties:
      name: { get_param: name }
      image: cirros
      flavor: { get_param: flavor_name }
      networks:
        - port: { get_resource: server1_port }

outputs:
  server1_private_ip:
    description: Private IP address of server1
    value: { get_attr: [ server1, first_address ] }
  server1_public_ip:
    description: Floating IP address of server1
    value: { get_attr: [ server1_floating_ip, floating_ip_address ] }`

    Tomcat = `heat_template_version: '2013-05-23'

parameters:
  name:
    type: string
    description: Name of VM

  net_name:
    type: string
    description: Name of network

  subnet_name:
    type: string
    description: Name of subnet

  flavor_name:
    type: string
    description: Name of the flavor

resources:
  server_security_group:
    type: OS::Neutron::SecurityGroup
    properties:
      description: Add security group rules for server
      name: tomcat
      rules:
        - remote_ip_prefix: 0.0.0.0/0
          protocol: tcp
          port_range_min: 8080
          port_range_max: 8080
        - remote_ip_prefix: 0.0.0.0/0
          protocol: icmp

  server1_port:
    type: OS::Neutron::Port
    properties:
      network_id: { get_param: net_name }
      security_groups: [ { get_resource: server_security_group }, default, ssh-ping ]
      fixed_ips:
        - subnet_id: { get_param: subnet_name }

  server1_floating_ip:
    type: OS::Neutron::FloatingIP
    properties:
      floating_network_id: provider
      port_id: { get_resource: server1_port }

  server1:
    type: OS::Nova::Server
    properties:
      name: { get_param: name }
      image: tomcat
      flavor: { get_param: flavor_name }
      networks:
        - port: { get_resource: server1_port }

outputs:
  server1_private_ip:
    description: Private IP address of server1
    value: { get_attr: [ server1, first_address ] }
  server1_public_ip:
    description: Floating IP address of server1
    value: { get_attr: [ server1_floating_ip, floating_ip_address ] }`
)