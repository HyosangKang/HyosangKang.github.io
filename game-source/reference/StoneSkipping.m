z = 50 ; %input ('x���� �ӵ��� �Է��ϼ��� : ') ;
image_size_x = 6.*z ; image_size_y = 2 ; % ��ǥ�� ũ��

theta = pi / 200 ; % ���� ������ ����
size_of_stone = 5 ; % ���� ũ��
mass = 1 ; % ���� ��
G = 10 ; % �߷� ���ӵ�
% ���� ��ġ ������ �Ʒ� �� �κ��̴�.
x = -image_size_x ; y = .5 ; % ���� �ʱ� ��ġ
vx =3*z ; vy = 0.5 ; % ���� �ʱ� �ӵ�
dt = 0.005 ; % ª�� �ð� ����

rho = 1 ; % ���� �е�
C1 = .1 ; % ��� ����
C2 = .3 ; % ���� ����


vx_data = vx ; vy_data = vy ;
t_data = 0 ;
data_index = 1 ;

% ���� �׸���
[ X Y ] = meshgrid ( linspace(-image_size_x,image_size_x,25), linspace(-image_size_x,image_size_x,25)) ; Z = zeros(size(X)) ;

% ���ɿ� �׸���
circle_theta = linspace ( 0, 2*pi, 100 ) ;
circle = 0 ; % ���ɿ� üũ ����
radius = 0 ; % ���ɿ� ������ ũ��
circle_x = 0 ; % ���ɿ� �߽� ��ġ
number_of_circle = 0 ; % ���ɿ� ����


while ( y  > -image_size_y )

    % line ( [ x , x + size_of_stone * cos(theta) ] , [ y, y + size_of_stone * sin(theta) ] , 'LineWidth' , 4) ;

    if ( y > 0 ) % ���� ���߿� �� ���� �� ��ӵ� ���
        vy = vy - G * dt ;

        if ( circle == 1 )
        circle = 0 ; % ������ �� �� ���ɿ� üũ �ʱ�ȭ
        end

    else % ���� ���鿡 ��� �� ��ӵ� ���

        if ( circle == 0 ) % ���� ���鿡 ����� ��
            number_of_circle = number_of_circle + 1 ; % ������� ��� ���ɿ��� ����
            circle_x(number_of_circle) = x  ; % ���ɿ��� �׷��� ��ġ�� ǥ��
            radius(number_of_circle) = 0 ; % ���ɿ��� �ʱ� ������
            circle = 1 ; % ���ǹ��� ���� ���� �� �ѹ� ���P�ǵ��� ��
        end

        if ( y + size_of_stone * sin(theta) > 0 ) % ���� ���鿡 �Ϻθ� ����� �� ���ӵ� ���
            Sim = min ( abs(y) , size_of_stone ) / sin(theta) ; % ���� ���鿡 ��� ����
            vvy = -G + 1/2 * rho / mass * ( vy^2 + vx^2 ) * Sim * ( C1 * cos(theta) - C2 * sin(theta) ) ; % ���������� ���ӵ�
            vvx = -1/2 * rho / mass * ( vy^2 + vx^2 ) * Sim * ( C1 * sin(theta) + C2 * cos(theta) ) ; % ��������� ���ӵ�

        else % ���� ������ ����� �� ���ӵ� ���
            vvx = -C2 * vx^2 ;
            vvy = -G + C2 * vy^2 ;
        end

        vx = vx + vvx * dt ;
        vy = vy + vvy * dt ;

    end

    x = x + vx * dt ;
    y = y + vy * dt ;

%     subplot (3,1,1); % ���� ����� ���ɿ� �׸���

    plot3 ( [ x , x + size_of_stone * cos(theta) ] , [0,0], [ y, y + size_of_stone * sin(theta) ] , 'LineWidth' , 4, 'Color', 'r') ; % �� �׸���
    hold on ;
    axis ( [ -image_size_x ,image_size_x , -image_size_x, image_size_x, -image_size_y , image_size_y ] ) ;

    if ( number_of_circle > 0 ) % ���ɿ�� �׸���
        for i = 1 : number_of_circle
            plot3( circle_x(i) + radius(i) * cos(circle_theta), radius(i) * sin(circle_theta), zeros(size(circle_theta)) ) ;
            radius(i) = radius(i) + 1.5 ; % 1.5�� ���ɿ��� Ŀ��
            %r = 1/(radius(i)+1) * sqrt((X-circle_x(i)).^2+Y.^2) ;
            %Z = Z - 0.5 * sin(r)./ exp(r) ;
        end
    end

    mesh( X, Y, Z ) ;
    Z = zeros ( size(X) ) ;

    xlabel ('distance') ;
    zlabel ('height');
    title ('Trajectory of stone') ;
    hold off;

    t_data = [ t_data, data_index*dt] ;
    vx_data = [ vx_data, vx ] ;
    vy_data = [ vy_data, vy ] ;
    data_index = data_index + 1 ;

%     subplot (3,1,2); % X ���� �ӵ� �׷���
%     plot ( t_data(:), vx_data(:)) ;
%     axis ( [ 0,10,0,8*z ] ) ;
%     grid on;
%
%     xlabel ('time') ;
%     ylabel ('velocity') ;
%     title ('Velocity of stone in X direction') ;
%
%     subplot(3,1,3) ; % Y ���� �ӵ� �׷���
%     plot ( t_data(:), vy_data(:)) ;
%     axis ( [ 0,10,-10,10] ) ;
%     grid on;
%
%     xlabel ('time') ;
%     ylabel ('velocity') ;
%     title ('Velocity of stone in Y direction') ;
%
    drawnow ;




end
